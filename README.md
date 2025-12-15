# conf2025badge-examples
<!-- pandoc -f markdown -t html5 -o README.html -c github.css README.md -->

RP2040チップを搭載した電子ネームタグ[conf2025badge](https://github.com/sago35/keyboards/blob/main/conf2025badge/build/build.md) 用に作成したTinygoのサンプルプログラム集です。アイデアの検証やテスト用に作成したものなので、実用性はありません。

## Hardware

Go Workshop Conference 2025 IN KOBE のワークショップ TinyGo Keeb Tour at GWCで、頒布された電子ネームタグ[conf2025badge](https://github.com/sago35/keyboards/blob/main/conf2025badge/build/build.md) を使用しました。  

![conf2025badge](photo/DSCN0350_800x600.jpg)

この電子ネームタグは、Raspberry Pi Pico と同じRP2040チップを搭載したマイコンボード[Seeed Studio XIAO RP2040](https://wiki.seeedstudio.com/XIAO-RP2040/)で作られています。  
ここで公開しているソースコードは、RP2040チップを使用しているマイコンボードであれば、多少改変すれば、利用できると思います。  

## Tinygo開発環境のインストール  

ここで、公開しているソースコードを利用するために必要なTinygoの開発環境を導入して下さい。  
ここでは、Windows11上での開発環境構築について解説します。他のOSについては、本家[tinygo](https://tinygo.org/)サイトの解説をお読み下さい。  

1. パッケージ管理ツールscoopのサイトを開き、導入スクリプトを入手して下さい。  

	[scoop](https://github.com/ScoopInstaller/Scoop)

2. Powershellを開いて、以下のスクリプトを実行して下さい。  

	\> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser  
	\> Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression  

3. 以下のコマンドを実行して、環境構築は終了です。

	\>scoop install go tinygo

4. 以下のコマンドを実行できれば、正常にインストールできています。  

	\>tinygo version  
	tinygo version 0.39.0 windows/amd64  
	(using go version go1.25.0 and LLVM version 19.1.2)

## Examples
<!-- 
### dispQRcode

電子ネームタグ[conf2025badge](https://github.com/sago35/keyboards/blob/main/conf2025badge/build/build.md) のOLEDディスプレイに、QRコードを表示します。  

* 解説 [./dispQRcode/README.md](./dispQRcode/README.md)  
* ソースコード [./dispQRcode/main.go](./dispQRcode/main.go)  

### QR おみくじ

電子ネームタグ[conf2025badge](https://github.com/sago35/keyboards/blob/main/conf2025badge/build/build.md) のOLEDディスプレイに、御神託をQRコードで表示する**おみくじ**です。  
スマホ等のQRコードリーダーで読み取って下さい。  

* 解説 [./QR-omikuji/README.md](./QR-omikuji/README.md)  
* ソースコード [./QR-omikuji/main.go](./QR-omikuji/main.go)  
-->

### タイマー

電子ネームタグ[conf2025badge](https://github.com/sago35/keyboards/blob/main/conf2025badge/build/build.md) を使用したシンプルなカウントダウンタイマーです。  
ロータリーエンコーダだけで操作できます。  

* 解説 [./timer/README.md](./timer/README.md)  
* ソースコード [./timer/main.go](./timer/main.go)  
