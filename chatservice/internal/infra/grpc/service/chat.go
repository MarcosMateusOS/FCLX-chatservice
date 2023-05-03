package service

import (
	"github.com/MarcosMateusOS/fclx/chatservice/internal/infra/grpc/pb"
	"github.com/MarcosMateusOS/fclx/chatservice/internal/usecase/chatcompleationstream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamUseCase chatcompleationstream.ChatCompletionUseCase
	ChatConfigStream            chatcompleationstream.ChatCompletionConfigInputDTO
	StreamChannel               chan chatcompleationstream.ChatCompletionOutputDTO
}

func NewChatService(
	chatCompletionStreamUseCase chatcompleationstream.ChatCompletionUseCase,
	chatConfigStream chatcompleationstream.ChatCompletionConfigInputDTO,
	streamChannel chan chatcompleationstream.ChatCompletionOutputDTO) *ChatService {

	return &ChatService{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		StreamChannel:               streamChannel,
	}

}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompleationstream.ChatCompletionConfigInputDTO{
		Model:                c.ChatConfigStream.Model,
		ModelMaxTokens:       c.ChatConfigStream.ModelMaxTokens,
		Temperature:          c.ChatConfigStream.Temperature,
		TopP:                 c.ChatConfigStream.TopP,
		N:                    c.ChatConfigStream.N,
		Stop:                 c.ChatConfigStream.Stop,
		MaxTokens:            c.ChatConfigStream.MaxTokens,
		InitialSystemMessage: c.ChatConfigStream.InitialSystemMessage,
	}

	input := chatcompleationstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()

	go func() {
		for msg := range c.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}

	}()

	_, err := c.ChatCompletionStreamUseCase.Execute(ctx, input)

	if err != nil {
		return err
	}

	return nil

}
