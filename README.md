# wafoon

## これなに

* AWS WAF の Web ACL の一覧を取得します
* AWS WAF の Web ACL のデフォルトアクションを変更出来ます (ALLOW 又は BLOCK)
* direnv 等と組み合わせて利用してください

## 使い方

### インストール

https://github.com/inokappa/wafoon/releases から環境に応じたバイナリをダウンロードしてください.

```
wget https://github.com/inokappa/wafoon/releases/download/v0.0.1/wafoon_darwin_amd64 -O ~/bin/wafoon
chmod +x ~/bin/wafoon
```

### ヘルプ

```sh
$ wafoon -h
Usage of wafoon:
  -aclid string
        Web ACL ID を指定.
  -allow
        Default Action を Allow に変更.
  -block
        Default Action を Block に変更
  -endpoint string
        AWS API のエンドポイントを指定.
  -profile string
        Profile 名を指定.
  -region string
        Region 名を指定. (default "ap-northeast-1")
  -version
        バージョンを出力.
```

### Web ACL 一覧取得

```sh
$ wafoon
+---------------------+--------------------------------------+---------------+
|        NAME         |               WEBACLID               | DEFAULTACTION |
+---------------------+--------------------------------------+---------------+
| dummy-inokara-waf   | xxxxxxxx-24a6-46ad-949c-zzzzzzzzzzzz | BLOCK         |
| dummy-inokara-waf   | aaaaaaaa-1ded-4ef1-bd00-bbbbbbbbbbbb | BLOCK         |
+---------------------+--------------------------------------+---------------+
```

### Web ACL デフォルトアクションの変更

```sh
$ wafoon -aclid=aaaaaaaa-1ded-4ef1-bd00-bbbbbbbbbbbb -allow
処理を続行しますか? (y/n): y
処理を続行します.
デフォルトアクションを変更しました.
+---------------------+--------------------------------------+---------------+
|        NAME         |               WEBACLID               | DEFAULTACTION |
+---------------------+--------------------------------------+---------------+
| dummy-inokara-waf   | xxxxxxxx-24a6-46ad-949c-zzzzzzzzzzzz | BLOCK         |
| dummy-inokara-waf   | aaaaaaaa-1ded-4ef1-bd00-bbbbbbbbbbbb | ALLOW         |
+---------------------+--------------------------------------+---------------+

$ wafoon -aclid=aaaaaaaa-1ded-4ef1-bd00-bbbbbbbbbbbb -block
処理を続行しますか? (y/n): y
処理を続行します.
デフォルトアクションを変更しました.
+---------------------+--------------------------------------+---------------+
|        NAME         |               WEBACLID               | DEFAULTACTION |
+---------------------+--------------------------------------+---------------+
| dummy-inokara-waf   | xxxxxxxx-24a6-46ad-949c-zzzzzzzzzzzz | BLOCK         |
| dummy-inokara-waf   | aaaaaaaa-1ded-4ef1-bd00-bbbbbbbbbbbb | BLOCK         |
+---------------------+--------------------------------------+---------------+
```