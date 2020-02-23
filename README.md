# goServerClient
ファイルをパケット単位に分割して転送するプログラム
(研究で使用するプログラムをGo言語に変更)

## ディレクトリ構成

- typefile
    構造体の定義

- retrieval, trasnfer
    送受信に用いるファイル

## 実行方法

- 送信者
    ```
    go run send.go
    ```

- 受信者
    ```
    go run recv.go
    ```