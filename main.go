package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

func setContents(screen tcell.Screen, x int, y int, str string, style tcell.Style) {
	for _, r := range str {
		screen.SetContent(x, y, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
}
func main() {
	filePath := "numbers.txt"
	buyLottery := NewBuyLottery(filePath)
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()
	screen.SetContent(0, 0, '_', nil, tcell.StyleDefault)

	quit := make(chan struct{})
	go func() {
		for {
			screen.Clear()
			setContents(screen, 33, 10, "1. 宝くじを購入する", tcell.StyleDefault)
			setContents(screen, 33, 12, "2. 購入した宝くじを確認する", tcell.StyleDefault)
			setContents(screen, 33, 14, "3. 宝くじの抽選を行う", tcell.StyleDefault)
			setContents(screen, 33, 16, "4. リセット", tcell.StyleDefault)
			setContents(screen, 33, 18, "5. 終了する", tcell.StyleDefault)

			setContents(screen, 33, 0, "入力：", tcell.StyleDefault)
			screen.Show()
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
					close(quit)
					return
				case tcell.KeyRune:
					input := ev.Rune()
					screen.SetContent(38, 0, ev.Rune(), nil, tcell.StyleDefault)
					screen.Show()
					es := screen.PollEvent()
					switch es := es.(type) {
					case *tcell.EventKey:
						switch es.Key() {
						case tcell.KeyBackspace:
							continue
						}
					}
					if input == '1' {
						screen.Clear()
						setContents(screen, 33, 0, "何枚購入しますか？", tcell.StyleDefault)
						screen.Show()
						eq := screen.PollEvent()
						switch eq := eq.(type) {
						case *tcell.EventKey:
							switch eq.Key() {
							case tcell.KeyBackspace:
								continue
							case tcell.KeyEscape, tcell.KeyEnter:
								close(quit)
								return
							case tcell.KeyRune:
								inpute, err := strconv.Atoi(string(eq.Rune()))
								inputes := strconv.Itoa(inpute)
								if err != nil {
									// エラー文
									//fmt.Println("無効な入力")
									continue
								}
								setContents(screen, 50, 0, inputes, tcell.StyleDefault)
								screen.Show()
								es := screen.PollEvent()
								switch es := es.(type) {
								case *tcell.EventKey:
									switch es.Key() {
									case tcell.KeyBackspace:
										continue
									}
								}
								for i := 0; i < inpute; i++ {
									nb := strconv.Itoa(i + 1)
									screen.Clear()
									setContents(screen, 33, 10, nb, tcell.StyleDefault)
									setContents(screen, 34, 10, "枚目", tcell.StyleDefault)
									setContents(screen, 33, 12, "1. 数字選択", tcell.StyleDefault)
									setContents(screen, 33, 14, "2. ランダム選択", tcell.StyleDefault)
									screen.Show()
									ex := screen.PollEvent()
									switch ex := ex.(type) {
									case *tcell.EventKey:
										switch ex.Key() {
										case tcell.KeyEscape, tcell.KeyEnter:
											close(quit)
											return
										case tcell.KeyBackspace:
											continue
										case tcell.KeyRune:
											inputs, err := strconv.Atoi(string(ex.Rune()))
											if err != nil {
												fmt.Println("無効な入力")
												continue
											}
											switch inputs {
											case 1:
												buyLottery.CBuyTicket(screen, tcell.StyleDefault)
												//buyLottery.BuyTicket(screen, tcell.StyleDefault)
											case 2:
												buyLottery.BuyTicket(screen, tcell.StyleDefault)
											default:
												setContents(screen, 33, 16, "無効な選択です", tcell.StyleDefault)
											}

										}
									}
								}
							}
						}

					}
					if input == '2' {
						l := 0
						w := 10
						var ns []string
						//fmt.Println("購入した宝くじ一覧:")
						screen.Clear()
						setContents(screen, l, w, "購入した宝くじ一覧:", tcell.StyleDefault)
						//setContents(screen, 33, 12, "2. 購入した宝くじを確認する", tcell.StyleDefault)
						//setContents(screen, 33, 14, "3. 宝くじの抽選を行う", tcell.StyleDefault)
						//setContents(screen, 33, 16, "4. リセット", tcell.StyleDefault)
						//setContents(screen, 33, 18, "5. 終了する", tcell.StyleDefault)
						//setContents(screen, 33, 0, "入力：", tcell.StyleDefault)
						//screen.Show()
						for _, lottery := range buyLottery.numbers {
							//fmt.Println(lottery.number)
							numbersStr := ""
							for _, num := range lottery.number {
								numbersStr += strconv.Itoa(num) + " "
							}
							ns = append(ns, numbersStr)
							//setContents(screen, 33, s+2, numbersStr, tcell.StyleDefault)
						}
						for _, m := range ns {
							w += 2
							if w == 32 {
								l += 30
								w = 12
							}
							setContents(screen, l, w, m, tcell.StyleDefault)

						}
						screen.Show()
						es := screen.PollEvent()
						switch es := es.(type) {
						case *tcell.EventKey:
							switch es.Key() {
							case tcell.KeyEscape, tcell.KeyEnter:
								continue
							}
						}
					}
					if input == '3' {
						screen.Clear()
						result := drawLottery()
						//検証用
						//nresult := []int{10, 11, 12, 13, 14, 15, 16}
						//result = nresult
						//

						//setContents(screen, 33, 10, fmt.Sprintf("数字を選択してください(1~37): %c", inputRune1), tcell.StyleDefault)
						setContents(screen, 0, 10, fmt.Sprintf("抽選結果: %v\n", result), tcell.StyleDefault)
						CheckNumber(buyLottery.numbers, result, screen, tcell.StyleDefault)
						screen.Show()
						es := screen.PollEvent()
						switch es := es.(type) {
						case *tcell.EventKey:
							switch es.Key() {
							case tcell.KeyEscape, tcell.KeyEnter:
								continue
							}
						}
						//CheckNumber(buyLottery.numbers, result, screen, tcell.StyleDefault)
					}
					if input == '4' {
						screen.Clear()
						Fclear(filePath)
						setContents(screen, 0, 10, "購入した宝くじをリセットしました。", tcell.StyleDefault)
						screen.Show()
						buyLottery = NewBuyLottery(filePath)
						es := screen.PollEvent()
						switch es := es.(type) {
						case *tcell.EventKey:
							switch es.Key() {
							case tcell.KeyBackspace:
								continue
							}
						}
					}
					if input == '5' {
						close(quit)
					}
					//fmt.Print(input)
					screen.SetContent(40, 0, input, nil, tcell.StyleDefault)
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}()
	<-quit
}
