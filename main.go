package main

import (
	"fmt"
	"log"
	"strconv"

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
	wstyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	rstyle := tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()
	screen.SetContent(0, 0, '_', nil, wstyle)

	lot := make(chan struct{})
	go func() {
		for {
			screen.Clear()
			setContents(screen, 33, 10, "1. 宝くじを購入する", wstyle)
			setContents(screen, 33, 12, "2. 購入した宝くじを確認する", wstyle)
			setContents(screen, 33, 14, "3. 宝くじの抽選を行う", wstyle)
			setContents(screen, 33, 16, "4. リセット", wstyle)
			setContents(screen, 33, 18, "5. 終了する", wstyle)
			setContents(screen, 33, 0, "入力：", wstyle)
			screen.Show()
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
					continue
				case tcell.KeyRune:
					input := ev.Rune()
					screen.SetContent(39, 0, ev.Rune(), nil, wstyle)
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
						setContents(screen, 33, 0, "何枚購入しますか(1～9):", wstyle)
						screen.Show()
						eq := screen.PollEvent()
						switch eq := eq.(type) {
						case *tcell.EventKey:
							switch eq.Key() {
							case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
								continue
							case tcell.KeyRune:
								screen.SetContent(57, 0, eq.Rune(), nil, wstyle)
								screen.Show()
								inpute, err := strconv.Atoi(string(eq.Rune()))
								if err != nil {
									setContents(screen, 33, 2, "無効な入力です", rstyle)
									setContents(screen, 33, 4, "適当な入力でメニュー画面に戻ります", rstyle)
									screen.Show()
									swith(screen)
								}
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
									setContents(screen, 33, 10, nb, wstyle)
									setContents(screen, 34, 10, "枚目:", wstyle)
									setContents(screen, 33, 12, "1. 数字選択", wstyle)
									setContents(screen, 33, 14, "2. ランダム選択", wstyle)
									screen.Show()
									ex := screen.PollEvent()
									switch ex := ex.(type) {
									case *tcell.EventKey:
										switch ex.Key() {
										case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyBackspace:
											continue
										case tcell.KeyRune:
											screen.SetContent(40, 10, ex.Rune(), nil, wstyle)
											screen.Show()
											inputs, err := strconv.Atoi(string(ex.Rune()))
											if err != nil {
												setContents(screen, 33, 16, "無効な選択です", rstyle)
												screen.Show()
												swith(screen)
											}
											switch inputs {
											case 1:
												buyLottery.CBuyTicket(screen)
											case 2:
												buyLottery.BuyTicket(screen)
											default:
												setContents(screen, 33, 16, "無効な選択です", rstyle)
												screen.Show()
												i--
												swith(screen)
											}

										}
									}
								}
							}
						}

					} else if input == '2' {
						l := 0
						w := 10
						var ns []string
						screen.Clear()
						setContents(screen, l, w, "購入した宝くじ一覧:", wstyle)
						setContents(screen, l, 34, "エンターキーでメニューに戻ります", wstyle)
						for _, lottery := range buyLottery.numbers {
							numbersStr := ""
							for _, num := range lottery.number {
								numbersStr += strconv.Itoa(num) + " "
							}
							ns = append(ns, numbersStr)
						}
						for _, m := range ns {
							w += 2
							if w == 32 {
								l += 30
								w = 12
							}
							setContents(screen, l, w, m, wstyle)

						}
						screen.Show()
						swith(screen)
					} else if input == '3' {
						screen.Clear()
						result := drawLottery()
						//検証用
						//nresult := []int{10, 11, 12, 13, 14, 15, 16}
						//result = nresult
						setContents(screen, 0, 10, fmt.Sprintf("抽選結果: %v\n", result), wstyle)
						CheckNumber(buyLottery.numbers, result, screen)
						setContents(screen, 0, 34, "エンターキーでメニューに戻ります", wstyle)
						screen.Show()
						es := screen.PollEvent()
						switch es := es.(type) {
						case *tcell.EventKey:
							switch es.Key() {
							case tcell.KeyEscape, tcell.KeyEnter:
								continue
							}
						}
					} else if input == '4' {
						screen.Clear()
						Fclear(filePath)
						setContents(screen, 0, 10, "購入した宝くじをリセットしました", wstyle)
						setContents(screen, 0, 12, "エンターキーでメニューに戻ります", wstyle)
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
					} else if input == '5' {
						close(lot)
					} else {
						setContents(screen, 33, 2, "無効な入力です", rstyle)
						screen.Show()
						swith(screen)
					}
				}
			}
		}
	}()
	<-lot
}
