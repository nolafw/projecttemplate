package vali

import "github.com/nolafw/validator/pkg/rule"

// ベースのルールは、projecttemplateに作っておく
var Required = rule.CreateRequired("必須です")
var MaxLength = rule.CreateMaxLengthStr("最大文字数10を超えています", 10)
