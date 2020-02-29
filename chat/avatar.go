package main

import (
	"errors"
)

// ErrNoAavatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラーです．
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません．")

// Avatarはユーザのプロフィール画像を表す型です．
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返します．
	// 問題が発生した場合にはエラーを返します．特にURLを取得できなかった場合には，ErrNoAvatarURLを返します．
	GetAvatarURL(c *client) (string, error)
}
