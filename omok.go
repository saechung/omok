package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// --- 1. 상수 및 타입 정의 ---

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

func (b *Board) PrintBoard() {
	fmt.Println("\n   00 01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18")
	for i := 0; i < BOARD_SIZE; i++ {
		fmt.Printf("%02d ", i)
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

// IsValid: 보드 범위 내인지
func IsValid(x, y int) bool {
	return x >= 0 && x < BOARD_SIZE && y >= 0 && y < BOARD_SIZE
}

// GetLineCount: (x, y)에 돌이 있다고 가정하고, (dx, dy) 방향(및 반대)의 총 개수
// (CheckWin과 CheckForbidden에서 사용)
func (b *Board) GetLineCount(x, y int, player int, dx, dy int) int {
	if !IsValid(x, y) || b[y][x] != player {
		return 0
	}

	count := 1 // (x, y)의 돌 포함
	// 정방향
	for i := 1; ; i++ {
		nx, ny := x+i*dx, y+i*dy
		if IsValid(nx, ny) && b[ny][nx] == player {
			count++
		} else {
			break
		}
	}
	// 역방향
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

// IsOpenThree: (x, y)에 놓인 돌(player)이 (dx, dy) 방향으로 '열린 3'인지 확인
// (가정: GetLineCount(x, y, player, dx, dy) == 3)
func (b *Board) IsOpenThree(x, y int, player int, dx, dy int) bool {
	// (x,y)는 player의 돌이 임시로 놓인 상태

	// 1. 정방향으로 끝의 *다음 칸* 찾기
	px, py := x, y
	for IsValid(px, py) && b[py][px] == player {
		px, py = px+dx, py+dy
	}
	// (px, py)는 연속된 돌의 (정방향) 끝 바로 다음 칸

	// 2. 역방향으로 끝의 *다음 칸* 찾기
	nx, ny := x, y
	for IsValid(nx, ny) && b[ny][nx] == player {
		nx, ny = nx-dx, ny-dy
	}
	// (nx, ny)는 연속된 돌의 (역방향) 끝 바로 다음 칸

	// 3. 양쪽이 모두 비어있는지(EMPTY) 확인
	positiveOpen := IsValid(px, py) && b[py][px] == EMPTY
	negativeOpen := IsValid(nx, ny) && b[ny][nx] == EMPTY

	return positiveOpen && negativeOpen
}

// --- 5. 승리 및 금수 조건 확인 ---

// CheckWin: 5목 승리(정확히 5개)를 확인
func (b *Board) CheckWin(lastX, lastY int, player int) bool {
	directions := [][2]int{
		{1, 0},  // 가로
		{0, 1},  // 세로
		{1, 1},  // 대각선
		{1, -1}, // 역대각선
	}

	for _, dir := range directions {
		count := b.GetLineCount(lastX, lastY, player, dir[0], dir[1])
		if count == 5 { // ★변경: >= 5 에서 == 5 로 (6목 금지)
			return true
		}
	}
	return false
}

// CheckForbidden: 흑돌(player)이 (x, y)에 두었을 때 금수(6목, 쌍삼)인지 확인
func (b *Board) CheckForbidden(x, y int, player int) (bool, string) {
	// (돌을 놓기 전에 검사하므로, 이 칸은 비어있어야 함)

	// 1. 임시로 돌을 놓음
	b[y][x] = player

	openThreeCount := 0
	directions := [][2]int{{1, 0}, {0, 1}, {1, 1}, {1, -1}}

	for _, dir := range directions {
		count := b.GetLineCount(x, y, player, dir[0], dir[1])

		// 2. 6목 검사
		if count > 5 {
			b[y][x] = EMPTY // 원상복구
			return true, "6목 이상 금지입니다."
		}

		// 3. 쌍삼 검사 (열린 3이 2개 이상)
		if count == 3 {
			if b.IsOpenThree(x, y, player, dir[0], dir[1]) {
				openThreeCount++
			}
		}
	}

	// 4. 임시로 놓은 돌 제거 (원상복구)
	b[y][x] = EMPTY

	// 5. 판정
	if openThreeCount >= 2 {
		return true, "쌍삼 금지입니다."
	}

	return false, "" // 금수 아님
}

// --- 3. 메인 게임 루프 ---

func main() {
	board := NewBoard()
	currentPlayer := BLACK
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Go 오목 게임 (19x19, 렌주룰 적용) ===")
	fmt.Println("좌표를 'x,y' (예: 9,9) 형식으로 입력하세요.")
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

		x, errX := strconv.Atoi(coords[0])
		y, errY := strconv.Atoi(coords[1])

		if errX != nil || errY != nil {
			fmt.Println("숫자로 좌표를 입력하세요.")
			continue
		}

		// --- ★ 로직 수정: 돌 놓기 및 검사 ---

		// 1. 기본 유효성 검사 (범위, 빈 칸)
		if !IsValid(x, y) {
			fmt.Println("잘못된 좌표입니다. (범위 초과)")
			continue
		}
		if board[y][x] != EMPTY {
			fmt.Println("이미 돌이 놓인 자리입니다.")
			continue
		}

		// 2. 금수 규칙 적용 (흑돌 차례에만)
		if currentPlayer == BLACK {
			isForbidden, reason := board.CheckForbidden(x, y, currentPlayer)
			if isForbidden {
				fmt.Printf("금수입니다: %s\n", reason)
				fmt.Println("다른 곳에 시도하세요.")
				continue // 플레이어 턴을 넘기지 않고 다시 입력받음
			}
		}

		// 3. (금수가 아니면) 돌 놓기
		board[y][x] = currentPlayer

		// 4. 승리 확인 (CheckWin은 5목만 확인하도록 수정됨)
		if board.CheckWin(x, y, currentPlayer) {
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