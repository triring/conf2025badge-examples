// tinygo build -o timer.uf2 -target=xiao-rp2040 -size short .
// tinygo flash -target=xiao-rp2040 -monitor -size short .

package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/encoders"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/tone"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	Waiting   int = 0 // 待機
	Countdown int = 1 // カウントダウン
	Reset     int = 2 // リセット
)

var (
	RemainingTime         time.Duration = 300  // 設定残り時間 5分間
	PreviousRemainingTime time.Duration = 300  // 設定残り時間 5分間
	MaxTime               time.Duration = 5940 // 最大時間
	MinTime               time.Duration = 0    // 最小時間
)

const (
	shortPressThresholdMs time.Duration = 100 // 短押しと判別する最大時間（ミリ秒）
	longPressThresholdMs  time.Duration = 500 // 長押しと判別する最小時間（ミリ秒）
)

// カラーユニバーサルデザイン(CUD) カラーセット
var (
	// Accent Colors アクセントカラー
	red      = color.RGBA{R: 0xFF, G: 0x4B, B: 0x0, A: 0xFF}  //  Red : 赤
	yellow   = color.RGBA{R: 0xFF, G: 0xF1, B: 0x0, A: 0xFF}  //  Yellow : 黄色
	green    = color.RGBA{R: 0x3, G: 0xAF, B: 0x7A, A: 0xFF}  //  Green : 緑
	blue     = color.RGBA{R: 0x0, G: 0x5A, B: 0xFF, A: 0xFF}  //  Blue : 青
	sky_blue = color.RGBA{R: 0x4D, G: 0xC4, B: 0xFF, A: 0xFF} //  Sky blue : 空色
	pink     = color.RGBA{R: 0xFF, G: 0x80, B: 0x82, A: 0xFF} //  Pink : ピンク
	orange   = color.RGBA{R: 0xF6, G: 0xAA, B: 0x0, A: 0xFF}  //  Orange : オレンジ
	purple   = color.RGBA{R: 0x99, G: 0x0, B: 0x99, A: 0xFF}  //  Purple : 紫
	brown    = color.RGBA{R: 0x80, G: 0x40, B: 0x0, A: 0xFF}  //  Brown : 茶色

	// Base Colors  ベースカラー
	light_pink         = color.RGBA{R: 0xFF, G: 0xCA, B: 0xBF, A: 0xFF} //  Light pink : 明るいピンク
	cream              = color.RGBA{R: 0xFF, G: 0xFF, B: 0x80, A: 0xFF} //  Cream : クリーム
	light_yellow_green = color.RGBA{R: 0xD8, G: 0xF2, B: 0x55, A: 0xFF} //  Light yellow-green : 明るい黄緑
	light_sky_blue     = color.RGBA{R: 0xBF, G: 0xE4, B: 0xFF, A: 0xFF} //  Light sky blue : 明るい空色
	beige              = color.RGBA{R: 0xFF, G: 0xCA, B: 0x80, A: 0xFF} //  Beige : ベージュ
	light_green        = color.RGBA{R: 0x77, G: 0xD9, B: 0xA8, A: 0xFF} //  Light green : 明るい緑
	light_purple       = color.RGBA{R: 0xC9, G: 0xAC, B: 0xE6, A: 0xFF} //  Light purple : 明るい紫

	// Achromatic Colors 無彩色
	white      = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} //  White  白
	light_gray = color.RGBA{R: 0xC8, G: 0xC8, B: 0xCB, A: 0xFF} //  Light gray  明るいグレー
	gray       = color.RGBA{R: 0x84, G: 0x91, B: 0x9E, A: 0xFF} //  Gray  グレー
	black      = color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xFF}    //  Black  黒
)

// Unix時間への相互変換（ミリ秒）

// システム時間からUnix時間への変換
func timeToUnixMilli(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

// Unix時間からシステム時間への変換
func unixMilliToTime(millis int64) time.Time {
	return time.Unix(0, millis*1000000)
}

func DispTime_oled(display ssd1306.Device, t time.Duration) {
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	display.ClearBuffer()
	DispTime := fmt.Sprintf("%02d:%02d", RemainingTime/60, RemainingTime%60)
	tinyfont.WriteLine(&display, &freemono.Bold9pt7b, 5, 24, "Timer", white)
	tinyfont.WriteLine(&display, &freemono.Bold18pt7b, 5, 60, DispTime, white)
	display.Display()
}

// tone.A7  7 octaves ド	2093.0 Hz
var mute tone.Note = 0 // 無音

// init関数は、パッケージがロードされたときに一度だけ実行されます。
/*
func init() {
	// ブザーが接続されたPinの設定と初期化
	bzrPin := machine.GPIO1
	pwm := machine.PWM0
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}
}
*/
func main() {
	// ブザーが接続されたPinの設定と初期化
	bzrPin := machine.GPIO1
	pwm := machine.PWM0
	speaker, err := tone.New(pwm, bzrPin)
	if err != nil {
		println("failed to configure PWM")
		return
	}

	// エンコーダの設定
	enc := encoders.NewQuadratureViaInterrupt(
		machine.GPIO3, // ROT_A1
		machine.GPIO4, // ROT_B1
	)

	enc.Configure(encoders.QuadratureConfig{
		Precision: 4,
	})

	rot_btn := machine.GPIO2 // ROT_BTN1
	rot_btn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	/* I2Cの初期設定 (Zero-kb02)
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 2.8 * machine.MHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})
	*/
	// I2Cの初期設定 (conf2025badge)
	machine.I2C1.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.GPIO6,
		SCL:       machine.GPIO7,
	})

	/* OLEDディスプレイの初期設定 (Zero-kb02)
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	*/

	// OLEDディスプレイの初期設定 (conf2025badge)
	display := ssd1306.NewI2C(machine.I2C1)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.SetRotation(drivers.Rotation0)
	display.ClearDisplay()

	var newValue int = 0
	var oldValue int = 0
	var state int = Waiting

	var timerStartTime time.Time
	var pressStartTime time.Time
	isButtonPressed := false
	buttonState := true

	pipoSound(speaker) // 起動音
	DispTime_oled(display, RemainingTime)
	for {
		// ボタンの状態を読み取る。
		// プルアップを使用しているので、ボタンが押されるとLOW（false）になる。
		buttonState = rot_btn.Get()
		if !buttonState && !isButtonPressed { // ボタンが押された瞬間
			pressStartTime = time.Now()
			isButtonPressed = true
			println("ボタンが押されました")
		} else if buttonState && isButtonPressed { // ボタンが離された瞬間
			pressDuration := time.Since(pressStartTime)
			isButtonPressed = false
			if pressDuration >= longPressThresholdMs*time.Millisecond {
				println(">>> 長押しを検出しました！ <<<")
				// reset
				RemainingTime = PreviousRemainingTime
				DispTime_oled(display, RemainingTime)
				pipoSound(speaker) // リセット音
			} else if pressDuration >= shortPressThresholdMs*time.Millisecond {
				println(">>> 短押しを検出しました！ <<<")
				if state == Waiting {
					if RemainingTime > 0 {
						state = Countdown
						PreviousRemainingTime = RemainingTime
						timerStartTime = time.Now()
						clickSound(speaker) // 受付・スタート音
					} else {
						println("リセットするか、ボタンを長押しして、計測する時間を再設定して下さい。")
						errorSound(speaker) // エラー音
					}
				} else {
					state = Waiting
					clickSound(speaker) // 受付・スタート音
				}
			} else {
				println("チャタリングまたは非常に短い押し込みとして無視します。")
			}
			println("ボタンが離されました。押下時間:", pressDuration)
		}
		switch state {
		case Waiting:
			if newValue = enc.Position(); newValue != oldValue {
				println("value: ", newValue)
				if (newValue - oldValue) > 0 {
					RemainingTime += time.Duration(newValue - oldValue)
					if RemainingTime > MaxTime {
						RemainingTime = MaxTime
					}
				} else {
					RemainingTime -= time.Duration(oldValue - newValue)
					if MinTime > RemainingTime {
						RemainingTime = MinTime
					}
				}
				oldValue = newValue
				DispTime_oled(display, RemainingTime)
			}
			break

		case Countdown:
			println(RemainingTime, time.Since(timerStartTime)/1000000000)
			RemainingTime = PreviousRemainingTime - time.Since(timerStartTime)/1000000000
			if RemainingTime >= 0 {
				DispTime_oled(display, RemainingTime)
			} else {
				state = Waiting
				// 終了音 3種類の中から選択する
				// 0 JIS0013 終了音 // 終了音（近）(4)
				// 1 thirori sound
				// 2 ultra man color timer
				choice := 1
				switch choice {
					case 0 : go playJISS0013end(speaker)
					case 1 : go playThirori(speaker)
					case 2 : go playColorTimer(speaker)
				}
			}
			break
		}
		// CPU使用率を抑えるため、短い間隔でポーリングする。
		time.Sleep(5 * time.Millisecond)
	}
}

// 受付・スタート音
func clickSound(speaker tone.Speaker) {
	speaker.SetNote(tone.A7) // 7 octaves ド	2093.0 Hz
	time.Sleep(time.Millisecond * 100)
	speaker.SetNote(mute)
}

// エラー音
func errorSound(speaker tone.Speaker) {
	speaker.SetNote(tone.A3) // 3 octaves ラ	 220.0 Hz
	time.Sleep(time.Millisecond * 1000)
	speaker.SetNote(mute)
}

// JIS S 0013:2022
// アクセシブルデザインー消費生活用製品の報知音
// 終了音（近）(4)
func playJISS0013end(speaker tone.Speaker) {
	var on1 time.Duration = 100
	var off1 time.Duration = 100
	var on2 time.Duration = 500
	var off2 time.Duration = 500

	for i := 0; i < 10; i++ {
		speaker.SetNote(tone.A7) // 7 octaves ド	2093.0 Hz
		time.Sleep(time.Millisecond * on1)
		speaker.SetNote(mute)
		time.Sleep(time.Millisecond * off1)
		speaker.SetNote(tone.A7) // 7 octaves ド	2093.0 Hz
		time.Sleep(time.Millisecond * on2)
		speaker.SetNote(mute)
		time.Sleep(time.Millisecond * off2)
	}
}

// 終了音 thirori
// ティロリ!!ティロリ!!ティロリ!!
// 某大手ハンバーガーチェーンで、ポテトが揚がったときに店内で流れるタイマー音
func playThirori(speaker tone.Speaker) {
//	var Repetitions int = 16             // 繰返しの回数
	var Repetitions int = 3             // 繰返しの回数
	// 楽曲のテンポは、125なので、8分音符
	for i := 0; i < Repetitions; i++ {
		speaker.SetNote(tone.G6)            // So
		time.Sleep(time.Millisecond * 240)  // eighth note
		speaker.SetNote(tone.E6)            // Mi
		time.Sleep(time.Millisecond * 240)  // eighth note
		speaker.SetNote(tone.G6)            // So
		time.Sleep(time.Millisecond * 240)  // eighth note
		speaker.SetNote(0)                  // 休符
		time.Sleep(time.Millisecond * 240)  // eighth note
	}
	speaker.SetNote(mute)
}

// 終了音 color timer
// ウルトラマンのカラータイマーの音
// A6 0.2秒、D6 0.4秒の繰り返し
// 
func playColorTimer(speaker tone.Speaker) {
	var Repetitions int = 20             // 繰返しの回数
	// 楽曲のテンポは、125なので、8分音符
	for i := 0; i < Repetitions; i++ {
		speaker.SetNote(tone.A6)            // A6 La 1749
		time.Sleep(time.Millisecond * 200)  // eighth note
		speaker.SetNote(tone.D6)            // D6 Re 1175 
		time.Sleep(time.Millisecond * 400)  // quarter note
	}
	speaker.SetNote(mute)
}

// 起動音、リセット音
// 国民機パソコン PC-9801シリーズの起動音	ピッポ!! を再現
func pipoSound(speaker tone.Speaker) {
	speaker.SetNote(tone.B6) // 6 octaves シ	PI 1980 Hz
	time.Sleep(time.Millisecond * 100)
	speaker.SetNote(mute)
	time.Sleep(time.Millisecond * 20)
	speaker.SetNote(tone.B5) // 5 octaves シ	PO 990 Hz
	time.Sleep(time.Millisecond * 100)
	speaker.SetNote(mute)
}
