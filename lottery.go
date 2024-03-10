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

// スペースで区切られた数字文字列を[]int に変換
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

func (ab *BuyLottery) BuyTicket(screen tcell.Screen) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	lottery := &Lottery{
		number: Rnumbers(),
	}
	if chs(lottery.number, ab) {
		ab.BuyTicket(screen)
	} else {
		ab.numbers = append(ab.numbers, lottery)
		ab.saveLottery()
		numbersStr := ""
		for _, num := range lottery.number {
			numbersStr += strconv.Itoa(num) + " "
		}
		setContents(screen, 33, 14, "ランダムで購入した宝くじ", style)
		setContents(screen, 33, 16, numbersStr, style)
		screen.Show()
		es := screen.PollEvent()
		switch es := es.(type) {
		case *tcell.EventKey:
			switch es.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
			}
		}
	}
}

func (ab *BuyLottery) CBuyTicket(screen tcell.Screen) {
	tnumber := make([]int, 7)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	fl := 0
	for i := 0; i < 7; i++ {
		screen.Clear()
		style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
		setContents(screen, 33, 10, "数字を選択してください(1~37): ", style)
		setContents(screen, 33, 12, "1桁の場合、01のように入力", style)
		setContents(screen, 19, 14, "選択した数字:", style)
		//現時点で選択されている数字一覧を表示
		h := 33
		for _, number := range tnumber {
			if number != 0 {
				pnumber := strconv.Itoa(number)
				setContents(screen, h, 14, pnumber, style)
				h += 3
			}
		}
		screen.Show()
		inputEvent1 := screen.PollEvent()
		switch inputEvent1 := inputEvent1.(type) {
		case *tcell.EventKey:
			switch inputEvent1.Key() {
			case tcell.KeyBackspace:
				i = 6
				fl = 1
				break
			case tcell.KeyRune:
				inputRune1 := inputEvent1.Rune()
				setContents(screen, 33, 10, fmt.Sprintf("数字を選択してください(1~37): %c", inputRune1), style)
				screen.Show()
				inputEvent2 := screen.PollEvent()
				switch inputEvent2 := inputEvent2.(type) {
				case *tcell.EventKey:
					switch inputEvent2.Key() {
					case tcell.KeyRune:
						inputRune2 := inputEvent2.Rune()
						setContents(screen, 33, 10, fmt.Sprintf("数字を選択してください(1~37): %c%c", inputRune1, inputRune2), style)
						screen.Show()
						swith(screen)
						input1, err1 := strconv.Atoi(string(inputRune1))
						input2, err2 := strconv.Atoi(string(inputRune2))
						if err1 != nil || err2 != nil {
							style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
							setContents(screen, 33, 16, "数字の変換に失敗しました", style)
							setContents(screen, 33, 18, "正しい数字を入力してください", style)
							screen.Show()
							swith(screen)
							i--
						} else {
							input := input1*10 + input2
							if contains(tnumber, input) || input < 1 || input > 37 {
								style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
								setContents(screen, 33, 16, "数字が重複、もしくは範囲外の数字です。 ", style)
								setContents(screen, 33, 18, "重複しない数字を選択してください", style)
								screen.Show()
								i--
								swith(screen)
							} else {
								tnumber[i] = input
								break
							}
						}
					}
				}
			}
		}
	}
	sort.Slice(tnumber, func(i, j int) bool {
		return tnumber[i] < tnumber[j]
	})
	lottery := &Lottery{
		number: tnumber,
	}
	if chs(tnumber, ab) && fl == 0 {
		style = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
		setContents(screen, 33, 16, "既に購入した宝くじです", style)
		setContents(screen, 33, 18, "別の宝くじを選択してください", style)
		screen.Show()
		swith(screen)
		ab.CBuyTicket(screen)
	} else {
		if !renumber(tnumber) {
			ab.numbers = append(ab.numbers, lottery)
			ab.saveLottery()
		}
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
		if v == num || 1 > num || num > 38 || num == 0 {
			return true
		}
	}
	return false
}

func renumber(slice []int) bool {
	for _, v := range slice {
		if v == 0 {
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

func CheckNumber(tickets []*Lottery, lotteryResult []int, screen tcell.Screen) {
	l := 0
	w := 10
	c := 0
	var q []int
	var ns []string
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
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
	for _, m := range ns {
		w += 2
		if w == 32 {
			l += 35
			w = 12
		}
		setContents(screen, l, w, m, style)
		cn, st := pn(q[c])
		setContents(screen, l+25, w, cn, st)
		c += 1
	}
}

func pn(n int) (string, tcell.Style) {
	switch n {
	case 0, 1, 2:
		return "はずれ", tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	case 3:
		return "5等", tcell.StyleDefault.Foreground(tcell.ColorBlue).Bold(true)
	case 4:
		return "4等", tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	case 5:
		return "3等", tcell.StyleDefault.Foreground(tcell.Color94).Bold(true)
	case 6:
		return "2等", tcell.StyleDefault.Foreground(tcell.ColorSilver).Bold(true)
	case 7:
		return "1等", tcell.StyleDefault.Foreground(tcell.ColorGold).Bold(true)
	default:
		return "", tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	}
}

func Fclear(filename string) {
	os.Remove(filename)
	os.Create(filename)
}

func swith(screen tcell.Screen) {
	for {
		es := screen.PollEvent()
		switch es := es.(type) {
		case *tcell.EventKey:
			switch es.Key() {
			case tcell.KeyEnter:
				return
			}
		}
	}
}
