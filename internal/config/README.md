ここの設定がデフォルトの設定になる。
各モジュール内のconfigの設定値がある場合は、そちらの設定に上書きされる

Parameter StoreとSecrets ManagerのARNも、ymlに記載する必要があるが、
`local.yml`はlocalstackのurlでいいが、他のは環境変数からインジェクトするか?
なんらかの方法で、直書きではなく、インジェクトする仕組みにしたい。
これも実行時の引数でしていできるか?

ymlで、どこから設定を読み込むかは、複数指定できるように、配列で設定できるようにする
環境変数からも読み込んで、secrets managerからも読み込む、というような使い方ができる。
ローカルでは、`.env`にしてしまうとかもできる。
```yml
# dotenv: .envファイルから読み込む
# env: 環境変数から読み込む
# ssmparam: parameter storeから読み込む
# secretsmanager: secrets managerから読み込む
var: [dotenv, env, ssmparam, secretsmanager]

```