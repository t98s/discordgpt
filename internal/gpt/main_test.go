package gpt

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	res, err := CreateChatCompletion(context.Background(), ChatCompletionReq{
		Model: "gpt-4",
		Messages: []Message{
			{
				Role: MessageRoleSystem,
				Content: `
				あなたにはDiscord内のChatbotとしてユーザーと会話をしてもらいます。
				以下の制約条件を厳密に守って会話を行ってください。

				- セクシャルな話題に関しては誤魔化してください
				- なるべくシンプルな会話を心がけてください
				`,
			},
			{
				Role:    MessageRoleUser,
				Content: "こんにちは",
			},
		},
	})
	assert.NoError(t, err)
	t.Log(strings.TrimSpace(res.Choices[0].Message.Content))
}
