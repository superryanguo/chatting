import sys
sys.path.append("../proto/")

from socket import AF_INET, SOCK_STREAM, socket
import chat_pb2 as chat

if __name__=="__main__":
    #ip = "localhost"
    ip = "127.0.0.1"
    HOST,PORT = ip,6789

    while(1):

        # ask question
        ask = chat.ChatAsk()
        sessionId = input("$Please input your sessionId: ")
        if sessionId.isdigit():
            ask.SessionId = sessionId
        else:
            ask.SessionId = "000"

        query = input("$Please input your query: ")
        if query == '1':
            ask.Query = "How are you?"
        else :
            ask.Query = query

        data=ask.SerializeToString()

        s = socket(AF_INET,SOCK_STREAM)
        s.connect((HOST,PORT))

        print("Client send OK: ",s.sendall(data))
        print("  SessionId: ", ask.SessionId)
        print("  Query: ", ask.Query)

        # get reply
        data_rcv=s.recv(1024)
        # print(b"Received:"+data_rcv)

        answer = chat.ChatAnswer()
        print("Client receive OK:",answer.ParseFromString(data_rcv))
        print("  SessionId: ", answer.SessionId)
        print("  Reply: ", answer.Reply)

        s.close()


