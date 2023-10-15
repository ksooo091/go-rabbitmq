package mailSender

import (
	"errors"
	"fmt"
	"net/smtp"
)

// SendMail to User
func SendMail(email, mailType, senderMail,senderPass string) {
	fmt.Println("전송시작")
	// 메일서버 로그인 정보 설정
	auth := smtp.PlainAuth("", senderMail, senderPass, "smtp.gmail.com")

	from := senderMail
	to := []string{email} // 복수 수신자 가능

	// 메시지 작성

	var headerSubject = ""
	var headerBlank = ""
	var body = ""
	switch mailType {
	case "F":
		headerSubject = "Subject: [Test]사진 생성에 문제가 생겼어요.\r\n"
		headerBlank = "\r\n"
		body = "요청하신 사진을 생성하던 중에 문제가생겼어요. 다시 신청해주세요.\r\n"
	case "S":
		headerSubject = "Subject: [Test]사진 생성이 완료 되었어요.\r\n"
		headerBlank = "\r\n"
		body = "요청하신 사진이 완성되었어요! 지금 접속하셔서 확인해보세요!\r\n"
	default:
		fmt.Println(errors.New("잘못된 메일타입"))
	}

	msg := []byte(headerSubject + headerBlank + body)

	if len(msg) != 0 {
		// 메일 보내기
		err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("%s에게 %s 형식의 메일 전송 완료\n", email, mailType)
		}
	}

}
