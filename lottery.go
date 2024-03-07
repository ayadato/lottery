package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// 宝くじ券
type Lottery struct {
	number []int
}

// 購入した宝くじ一覧
type BuyLottery struct {
	file    string
	numbers []*Lottery
}

func NewBuyLottery(file string) *BuyLottery {
	ab := &BuyLottery{
		file: file,
	}

	ab.readLottery()
	return ab
}

// Parse関数: スペースで区切られた数字文字列を[]int に変換
func Parse(input string) []int {
	var numbers []int
	for _, strNum := range strings.Fields(input) {
		num, err := strconv.Atoi(strNum)
		if err != nil {
			fmt.Fprintln(os.Stderr, "エラー：", err)
			os.Exit(1)
		}
		numbers = append(numbers, num)
	}
	return numbers
}

func (ab *BuyLottery) readLottery() {
	f, err := os.Open(ab.file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー：", err)
		os.Exit(1)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		numbers := Parse(s.Text())
		lottery := &Lottery{
			number: numbers,
		}

		ab.numbers = append(ab.numbers, lottery)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "エラー：", err)
		os.Exit(1)
	}
}

func (ab *BuyLottery) BuyTicket(screen tcell.Screen, style tcell.Style) {
	lottery := &Lottery{
		number: Rnumbers(),
	}
	if chs(lottery.number, ab) {
		ab.BuyTicket(screen, style)
	} else {
		ab.numbers = append(ab.numbers, lottery)
		ab.saveLottery()
		//fmt.Println("ランダムで購入した宝くじ:")
		//fmt.Println(lottery.number)

		numbersStr := ""
		for _, num := range lottery.number {
			numbersStr += strconv.Itoa(num) + " "
		}

		setContents(screen, 33, 14, "ランダムで購入した宝くじ", tcell.StyleDefault)
		setContents(screen, 33, 16, numbersStr, tcell.StyleDefault)
		screen.Show()
		es := screen.PollEvent()
		switch es := es.(type) {
		case *tcell.EventKey:
			switch es.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
			}
		}
		//time.Sleep(3 * time.Second)
	}
}

func (ab *BuyLottery) CBuyTicket(screen tcell.Screen, style tcell.Style) {
	numbers := make([]int, 7)
	for i := 0; i < 7; i++ {
		screen.Clear()
		setContents(screen, 33, 10, "数字を選択してください(1~37): ", tcell.StyleDefault)
		setContents(screen, 33, 12, "1桁の場合、01のように入力", tcell.StyleDefault)
		screen.Show()

		// 1桁目の数字を取得
		inputEvent1 := screen.PollEvent()
		switch inputEvent1 := inputEvent1.(type) {
		case *tcell.EventKey:
			switch inputEvent1.Key() {
			case tcell.KeyRune:
				inputRune1 := inputEvent1.Rune()

				// 1桁の数字を表示
				setContents(screen, 33, 10, fmt.Sprintf("数字を選択してください(1~37): %c", inputRune1), tcell.StyleDefault)
				screen.Show()

				//time.Sleep(1 * time.Second)

				// 2桁目の数字を取得
				inputEvent2 := screen.PollEvent()
				switch inputEvent2 := inputEvent2.(type) {
				case *tcell.EventKey:
					switch inputEvent2.Key() {
					case tcell.KeyRune:
						inputRune2 := inputEvent2.Rune()
						setContents(screen, 33, 10, fmt.Sprintf("数字を選択してください(1~37): %c%c", inputRune1, inputRune2), tcell.StyleDefault)
						screen.Show()
						es := screen.PollEvent()
						switch es := es.(type) {
						case *tcell.EventKey:
							switch es.Key() {
							case tcell.KeyEscape, tcell.KeyEnter:
							}
						}

						// 1桁の数字の場合、そのまま処理
						if inputRune2 >= '0' && inputRune2 <= '9' {
							input, err := strconv.Atoi(string(inputRune1) + string(inputRune2))
							if err != nil {
								setContents(screen, 33, 10, "数字の変換に失敗しました", tcell.StyleDefault)
								setContents(screen, 33, 12, "重複しない数字を選択してください", tcell.StyleDefault)
								screen.Show()
								continue
							}

							// 入力が範囲外の場合や重複している場合の処理
							for contains(numbers, input) || input < 1 || input > 37 {
								setContents(screen, 33, 10, "数字が重複、もしくは範囲外の数字です。 ", tcell.StyleDefault)
								setContents(screen, 33, 12, "重複しない数字を選択してください", tcell.StyleDefault)
								screen.Show()

								inputEvent := screen.PollEvent()
								switch inputEvent := inputEvent.(type) {
								case *tcell.EventKey:
									switch inputEvent.Key() {
									case tcell.KeyRune:
										inputRune2 := inputEvent.Rune()

										// 1桁の数字の場合、そのまま処理
										if inputRune2 >= '0' && inputRune2 <= '9' {
											input, err = strconv.Atoi(string(inputRune1) + string(inputRune2))
											if err != nil {
												setContents(screen, 33, 10, "数字の変換に失敗しました", tcell.StyleDefault)
												setContents(screen, 33, 12, "重複しない数字を選択してください", tcell.StyleDefault)
												screen.Show()
												continue
											}
										}
									}
								}
							}
							numbers[i] = input
							break
						}
					}
				}
			}
		}
	}
	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i] < numbers[j]
	})
	lottery := &Lottery{
		number: numbers,
	}
	if chs(numbers, ab) {
		setContents(screen, 33, 14, "既に購入した宝くじです。", tcell.StyleDefault)
		setContents(screen, 33, 16, "別の宝くじを選択してください。", tcell.StyleDefault)
		screen.Show()
		es := screen.PollEvent()
		switch es := es.(type) {
		case *tcell.EventKey:
			switch es.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
			}
		}
		//time.Sleep(3 * time.Second)
		//fmt.Println("既に購入した宝くじです。")
		//fmt.Println("別の宝くじを選択してください。")
		ab.CBuyTicket(screen, style)
	} else {
		ab.numbers = append(ab.numbers, lottery)
		ab.saveLottery()
	}
}

func (ab *BuyLottery) saveLottery() {
	f, err := os.Create(ab.file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー：", err)
		os.Exit(1)
	}
	defer f.Close()

	for _, lottery := range ab.numbers {
		for _, num := range lottery.number {
			_, err := fmt.Fprintf(f, "%d ", num)
			if err != nil {
				fmt.Fprintln(os.Stderr, "エラー：", err)
				os.Exit(1)
			}
		}
		_, _ = fmt.Fprintln(f)
	}
}

func drawLottery() []int {
	return Rnumbers()
}

func Rnumbers() []int {
	numbers := make([]int, 7)

	for i := 0; i < 7; i++ {
		// 既に取得した数字と重複しないように生成
		num := rand.Intn(37) + 1
		for contains(numbers, num) {
			num = rand.Intn(37) + 1
		}
		numbers[i] = num
	}
	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i] < numbers[j]
	})

	return numbers
}

func contains(slice []int, num int) bool {
	for _, v := range slice {
		if v == num || 0 > num || num > 38 {
			return true
		}
	}
	return false
}

func chs(slice []int, ab *BuyLottery) bool {
	for _, a := range ab.numbers {
		match := true
		for _, v := range slice {
			if !contains(a.number, v) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func CheckNumber(tickets []*Lottery, lotteryResult []int, screen tcell.Screen, style tcell.Style) {
	//screen.Clear()
	l := 0
	w := 10
	c := 0
	var q []int
	var ns []string
	//fmt.Println("当選した宝くじ番号:")
	for _, ticket := range tickets {
		q = append(q, compare(ticket.number, lotteryResult))
		//fmt.Println(ticket.number)
		//pn(q)
		numbersStr := ""
		for _, num := range ticket.number {
			numbersStr += strconv.Itoa(num) + " "
		}
		ns = append(ns, numbersStr)
	}
	//setContents(screen, 33, s+2, numbersStr, tcell.StyleDefault)
	for _, m := range ns {
		w += 2
		if w == 32 {
			l += 35
			w = 12
		}
		setContents(screen, l, w, m, tcell.StyleDefault)
		setContents(screen, l+25, w, pn(q[c]), tcell.StyleDefault)
		c += 1
	}
	screen.Show()
	es := screen.PollEvent()
	switch es := es.(type) {
	case *tcell.EventKey:
		switch es.Key() {
		case tcell.KeyEscape, tcell.KeyEnter:
		}
	}
}

func compare(numbers1, numbers2 []int) int {
	n := 0
	for i := range numbers1 {
		for j := range numbers2 {
			if numbers1[i] == numbers2[j] {
				n++
			}
		}
	}
	return n
}

func pn(n int) string {
	switch n {
	case 0, 1, 2:
		return "はずれ"
	case 3:
		return "5等"
	case 4:
		return "4等"
	case 5:
		return "3等"
	case 6:
		return "2等"
	case 7:
		return "1等"
	default:
		return ""
	}
}

func Fclear(filename string) {
	// ファイルを開く（存在しなければ新しく作成）
	os.Remove(filename)
	os.Create(filename)
}
