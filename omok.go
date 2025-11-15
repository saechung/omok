package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// --- 1. 상수 및 타입 정의 ---
// (변경 없음)
const (
	BOARD_SIZE = 19
	EMPTY      = 0
	BLACK      = 1
	WHITE      = 2
)

type Board [BOARD_SIZE][BOARD_SIZE]int

// --- 2. 보드 초기화 및 출력 ---

func NewBoard() *Board {
	var board Board
	return &board
}

// PrintBoard: 1-19 기반 좌표로 출력
func (b *Board) PrintBoard() {
	// ★변경: X좌표 헤더 (01 ~ 19)
	fmt.Print("\n   ")
	for i := 1; i <= BOARD_SIZE; i++ {
		fmt.Printf("%02d ", i)
	}
	fmt.Println()

	for i := 0; i < BOARD_SIZE; i++ {
		fmt.Printf("%02d ", i+1) // ★변경: Y좌표 (01 ~ 19)
		for j := 0; j < BOARD_SIZE; j++ {
			switch b[i][j] {
			case EMPTY:
				fmt.Print("┼  ")
			case BLACK:
				fmt.Print("●  ")
			case WHITE:
				fmt.Print("○  ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// --- 4. 유효성 검사 및 핵심 로직 헬퍼 ---
// (내부 로직은 0-18 기반이므로 변경 없음)

// IsValid: 보드 범위 내인지 (0 ~ 18 기준)
func IsValid(x, y int) bool {
	return x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE
}

// GetLineCount: (0-18 기준)
func (b *Board) GetLineCount(x, y int, player int, dx, dy int) int {
	if !IsValid(x, y) || b[y][x] != player {
		return 0
	}
	count := 1
	for i := 1; ; i++ {
		nx, ny := x+i*dx, y+i*dy
		if IsValid(nx, ny) && b[ny][nx] == player {
			count++
		} else {
			break
		}
	}
	for i := 1; ; i++ {
		nx, ny := x-i*dx, y-i*dy
		if IsValid(nx, ny) && b[ny][nx] == player {
			count++
		} else {
			break
		}
	}
	return count
}

// IsOpenThree: (0-18 기준)
func (b *Board) IsOpenThree(x, y int, player int, dx, dy int) bool {
	px, py := x, y
	for IsValid(px, py) && b[py][px] == player {
		px, py = px+dx, py+dy
	}
	nx, ny := x, y
	for IsValid(nx, ny) && b[ny][nx] == player {
		nx, ny = nx-dx, ny-dy
	}
	positiveOpen := IsValid(px, py) && b[py][px] == EMPTY
	negativeOpen := IsValid(nx, ny) && b[ny][nx] == EMPTY
	return positiveOpen && negativeOpen
}

// --- 5. 승리 및 금수 조건 확인 ---
// (내부 로직은 0-18 기반이므로 변경 없음)

// CheckWin: (0-18 기준)
func (b *Board) CheckWin(lastX, lastY int, player int) bool {
	directions := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
	for _, dir := range directions {
		count := b.GetLineCount(lastX, lastY, player, dir[0], dir[1])
		if count == 5 {
			return true
		}
	}
	return false
}

// CheckForbidden: (0-18 기준)
func (b *Board) CheckForbidden(x, y int, player int) (bool, string) {
	b[y][x] = player
	openThreeCount := 0
	directions := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		count := b.GetLineCount(x, y, player, dir[0], dir[1])
		if count > 5 {
			b[y][x] = EMPTY
			return true, "6목 이상 금지입니다."
		}
		if count == 3 {
			if b.IsOpenThree(x, y, player, dir[0], dir[1]) {
				openThreeCount++
			}
		}
	}
	b[y][x] = EMPTY
	if openThreeCount >= 2 {
		return true, "쌍삼 금지입니다."
	}
	return false, ""
}

// --- 3. 메인 게임 루프 ---

func main() {
	board := NewBoard()
	currentPlayer := BLACK
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Go 오목 게임 (19x19, 렌주룰 적용) ===")
	// ★변경: 안내 메시지 (1-19 기준)
	fmt.Println("!! 좌표를 'x,y' (예: 10,10) 형식으로 입력하세요. (1~19 범위) !!")
	fmt.Println("!! 흑돌은 6목, 쌍삼이 금지됩니다.")

	for {
		board.PrintBoard()

		var playerSymbol string
		var playerName string
		if currentPlayer == BLACK {
			playerSymbol = "●"
			playerName = "흑"
		} else {
			playerSymbol = "○"
			playerName = "백"
		}
		fmt.Printf("플레이어 %s (%s)의 차례입니다. (x,y): ", playerName, playerSymbol)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("입력 오류:", err)
			continue
		}

		coords := strings.Split(strings.TrimSpace(input), ",")
		if len(coords) != 2 {
			fmt.Println("잘못된 형식입니다. 'x,y' 형식으로 입력하세요.")
			continue
		}

		// 사용자가 입력한 1-19 좌표
		x, errX := strconv.Atoi(coords[0])
		y, errY := strconv.Atoi(coords[1])

		if errX != nil || errY != nil {
			fmt.Println("숫자로 좌표를 입력하세요.")
			continue
		}

		// ★★★ 핵심 변경 ★★★
		// 사용자가 입력한 1-19 좌표를 내부 0-18 좌표로 변환
		internalX := x - 1
		internalY := y - 1

		// --- ★ 로직 수정: 모든 검사와 처리는 internal 좌표계 사용 ---

		// 1. 기본 유효성 검사 (범위, 빈 칸)
		if !IsValid(internalX, internalY) { // ★변경
			fmt.Println("잘못된 좌표입니다. (1~19 범위 내로 입력하세요)")
			continue
		}
		if board[internalY][internalX] != EMPTY { // ★변경
			fmt.Println("이미 돌이 놓인 자리입니다.")
			continue
		}

		// 2. 금수 규칙 적용 (흑돌 차례에만)
		if currentPlayer == BLACK {
			isForbidden, reason := board.CheckForbidden(internalX, internalY, currentPlayer) // ★변경
			if isForbidden {
				fmt.Printf("금수입니다: %s\n", reason)
				fmt.Println("다른 곳에 시도하세요.")
				continue
			}
		}

		// 3. (금수가 아니면) 돌 놓기
		board[internalY][internalX] = currentPlayer // ★변경

		// 4. 승리 확인
		if board.CheckWin(internalX, internalY, currentPlayer) { // ★변경
			board.PrintBoard()
			fmt.Printf("축하합니다! 플레이어 %s (%s)의 승리입니다!\n", playerName, playerSymbol)
			break
		}

		// 5. 플레이어 교체
		if currentPlayer == BLACK {
			currentPlayer = WHITE
		} else {
			currentPlayer = BLACK
		}
	}
}