package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	proto "lab1_test"
	"math"
	_ "math"
	"net"
	"strings"
)

// Client - состояние клиента.
type Client struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений

}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			err := client.logger.Error("cannot decode message", "reason", err)
			if err != nil {
				return
			}
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "get":
		errorMsg := ""

		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {

			var a proto.TrigValue
			err := json.Unmarshal(*req.Data, &a)
			if err != nil {
				fmt.Println(err.Error())
				return false
			}
			if err == nil {

				val, err := getTrigValue(a.Fucn, a.Value)

				if err == nil && !math.IsNaN(val) {
					client.respond("result", &val)
					break
				} else if err != nil {
					errorMsg = "No such function\nUnable functions: sin, cos, tan, cot, asin, acos, atan, acot"
				} else if math.IsNaN(val) {
					errorMsg = fmt.Sprint(a.Fucn) + " undefined for this value"
				}
			}

		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.respond("failed", errorMsg)
		}
	default:
		err := client.logger.Error("unknown command")
		if err != nil {
			return false
		}
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	err := client.enc.Encode(&proto.Response{Status: status, Data: &raw})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func getTrigValue(function string, value float64) (float64, any) {

	function = strings.TrimSpace(function)
	//var rad float64 = 57.29578
	var rad float64 = 1
	switch function {
	case "sin":
		return math.Sin(value / rad), recover()
	case "cos":
		return math.Cos(value / rad), recover()
	case "tan":
		return math.Tan(value / rad), recover()
	case "cot":
		return 1 / math.Tan(value/rad), recover()
	case "asin":
		return math.Asin(value / rad), recover()
	case "acos":
		return math.Acos(value / rad), recover()
	case "atan":
		return math.Atan(value / rad), recover()
	case "acot":
		return math.Atan(1 / (value / 57.2956)), recover()
	default:
		return 0, fmt.Errorf("no such function")
	}
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string

	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)

		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {

					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.

					go NewClient(conn).serve()

				}
			}
		}
	}
}
