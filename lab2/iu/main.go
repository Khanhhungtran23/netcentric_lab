package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	// Library constants
	maxSeats       = 30
	totalStudents  = 100
	minStudyHours  = 1
	maxStudyHours  = 4
	simulationHour = time.Second // 1 second = 1 hour in simulation
)

// Student represents a library visitor
type Student struct {
	ID          int
	StudyHours  int
	ArrivalTime int
	LeaveTime   int
}

// Library represents the library with seats and operations
type Library struct {
	seats      int
	maxSeats   int
	waitingQ   []*Student
	studyingQ  []*Student
	hour       int
	totalHours int
	mu         sync.Mutex
	wg         sync.WaitGroup
}

// NewLibrary creates a new library instance
func NewLibrary(maxSeats int) *Library {
	return &Library{
		seats:      maxSeats,
		maxSeats:   maxSeats,
		waitingQ:   make([]*Student, 0, totalStudents),
		studyingQ:  make([]*Student, 0, maxSeats),
		hour:       0,
		totalHours: 0,
		mu:         sync.Mutex{},
	}
}

// EnterLibrary handles student entry to the library
func (l *Library) EnterLibrary(s *Student) {
	l.mu.Lock()
	s.ArrivalTime = l.hour

	if l.seats > 0 {
		// Seat available, student can enter
		l.seats--
		l.studyingQ = append(l.studyingQ, s)
		fmt.Printf("Time %d: Student %d starts reading at the lib\n", l.hour, s.ID)
		l.mu.Unlock() // Unlock before sleeping

		// Chỉ sleep và leave nếu student được vào thư viện
		time.Sleep(time.Duration(s.StudyHours) * simulationHour)
		l.LeaveLibrary(s)
	} else {
		// No seats available, student must wait
		l.waitingQ = append(l.waitingQ, s)
		fmt.Printf("Time %d: Student %d is waiting\n", l.hour, s.ID)
		l.mu.Unlock() // Unlock và không block thread với time.Sleep
		// Student đang chờ sẽ được xử lý trong LeaveLibrary khi có chỗ trống
	}
}

// LeaveLibrary handles student departure from the library
func (l *Library) LeaveLibrary(s *Student) {
	l.mu.Lock()

	foundStudent := false
	// Find and remove the student from studying queue
	for i, student := range l.studyingQ {
		if student.ID == s.ID {
			foundStudent = true
			// Remove student from studying queue
			l.studyingQ = append(l.studyingQ[:i], l.studyingQ[i+1:]...)
			l.seats++
			s.LeaveTime = l.hour
			fmt.Printf("Time %d: Student %d is leaving. Spent %d hours reading\n",
				l.hour, s.ID, s.StudyHours)

			// Check if any student is waiting
			if len(l.waitingQ) > 0 {
				// Get the first waiting student
				waitingStudent := l.waitingQ[0]
				l.waitingQ = l.waitingQ[1:]

				// Student enters the library
				l.seats--
				waitingStudent.ArrivalTime = l.hour
				l.studyingQ = append(l.studyingQ, waitingStudent)

				fmt.Printf("Time %d: Student %d stops waiting and starts reading\n",
					l.hour, waitingStudent.ID)

				// Tạo một bản sao của student và mutex unlock trước khi tạo goroutine mới
				studentCopy := waitingStudent
				l.wg.Add(1)
				l.mu.Unlock() // Mở khóa trước khi gọi goroutine mới

				go func(ws *Student) {
					defer l.wg.Done()
					time.Sleep(time.Duration(ws.StudyHours) * simulationHour)
					l.LeaveLibrary(ws)
				}(studentCopy)

				return // Thoát khỏi hàm sau khi đã unlock mutex
			}
			break
		}
	}

	l.mu.Unlock() // Đảm bảo mutex luôn được unlock

	if foundStudent {
		l.wg.Done() // Chỉ gọi Done() nếu tìm thấy và xử lý student
	}
}

// SimulateDay runs the library simulation for a day
func (l *Library) SimulateDay() {
	// Seed the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Create students with random study hours
	students := make([]*Student, totalStudents)
	for i := 0; i < totalStudents; i++ {
		students[i] = &Student{
			ID:         i + 1,
			StudyHours: r.Intn(maxStudyHours-minStudyHours+1) + minStudyHours,
		}
	}

	// Randomize student arrival order
	rand.Shuffle(len(students), func(i, j int) {
		students[i], students[j] = students[j], students[i]
	})

	// Tạo kênh để báo hiệu khi mô phỏng kết thúc
	done := make(chan bool)

	// Start the clock
	clockDone := make(chan bool) // Kênh để dừng goroutine đồng hồ
	go func() {
		for {
			select {
			case <-clockDone:
				return
			default:
				time.Sleep(simulationHour)
				l.mu.Lock()
				l.hour++
				l.mu.Unlock()
			}
		}
	}()

	// Send students to the library trong một goroutine riêng
	go func() {
		for _, student := range students {
			l.wg.Add(1)
			go l.EnterLibrary(student)

			// Small delay between student arrivals to simulate realistic flow
			time.Sleep(simulationHour / 10)
		}

		// Wait for all students to finish
		l.wg.Wait()

		// Báo hiệu mô phỏng đã hoàn thành
		done <- true
	}()

	// Đợi mô phỏng hoàn thành
	<-done

	// Dừng đồng hồ
	clockDone <- true

	// Calculate the total hours library was open
	maxLeaveTime := 0
	l.mu.Lock() // Lock khi đọc dữ liệu students
	for _, student := range students {
		if student.LeaveTime > maxLeaveTime {
			maxLeaveTime = student.LeaveTime
		}
	}
	l.mu.Unlock()

	l.totalHours = maxLeaveTime
	fmt.Printf("\nTime %d: No more students. Let's call it a day\n", l.totalHours)
	fmt.Printf("\nThe library needed to be open for %d hours to serve all %d students.\n",
		l.totalHours, totalStudents)
}

func main() {
	fmt.Println("International University Library Simulation")
	fmt.Printf("- Library capacity: %d seats\n", maxSeats)
	fmt.Printf("- Total students: %d\n", totalStudents)
	fmt.Printf("- Study time range: %d-%d hours\n\n", minStudyHours, maxStudyHours)

	library := NewLibrary(maxSeats)
	library.SimulateDay()
}
