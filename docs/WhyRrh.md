---
title: Why Rrh?
---

# Git repository manager

Gitリポジトリを管理するツールに[ghq](https://github.com/motemen/ghq)がある．
しかし，ghqを使い始めるにあたり次の事柄が問題となり，使い始められなかった．
なお，前提として，私のホームディレクトリにはすでに大量のgitリポジトリ（約350）がある．

* ghq の管理下に置くには clone しないといけない．
  私の環境で，全てのリポジトリを clone し直すのは現実的ではない．
* gitリポジトリの置き場所が特定の箇所に決められている．
  もちろん，複数箇所も対応しているようであるが，いくつかのディレクトリに分けて管理している環境では，置き場所を切り替えるのが面倒臭い．
  置き場所は自分で決めたい．

また，複数のリポジトリを同時並行で編集する場合もあるため，次のような機能があれば良いなと妄想する．

* どのリポジトリが最新で，どれが更新していないのか混乱する場合がある．
  そのため，複数の git リポジトリの状態（リモート，ローカル，インデックスの最終更新日時）が分かるようになれば良い．
* 複数リポジトリの内容を fetch したい．pull は内容を確認してからにしたいので，置いておく．

ということで作成してみたのが，Rrh である．