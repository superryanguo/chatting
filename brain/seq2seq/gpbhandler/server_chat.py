import sys
#sys.path.append("../proto/")

from socketserver import BaseRequestHandler, TCPServer
import chat_pb2 as chat
import time

chat_data = {}

class GPBHandler(BaseRequestHandler):
    def handle(self):

        # receive client data
        data=self.request.recv(1024)
        # print(b"Received:"+data)

        ask = chat.ChatAsk()
        print(time.strftime("%Y-%m-%d-%H:%M:%S", time.localtime()), "server receive OK:",ask.ParseFromString(data))
        print("  SessionId: ", ask.SessionId)
        print("  Query: ", ask.Query)

        # reply
        answer = chat.ChatAnswer()
        if ask.SessionId.isdigit():
            answer.SessionId = ask.SessionId
        else:
            answer.SessionId = "?"

        if ask.Query == "How are you?" :
            answer.Reply =  "I'm Fine. Hope you have a good day!"
        else :
            answer.Reply =  "Hello."

        print(time.strftime("%Y-%m-%d-%H:%M:%S", time.localtime()), "server reply ok:",self.request.sendall(answer.SerializeToString()))
        print("  SessionId: ", answer.SessionId)
        print("  Reply: ", answer.Reply)

        # save 'ask' and 'answer'word_freq.
        if chat_data.get(ask.SessionId) is not None:
            chat_data[ask.SessionId] = chat_data.get(ask.SessionId)+[("Ask:"+ ask.Query, "Ans:"+ answer.Reply)]
        else :
            chat_data[ask.SessionId] = [("Ask:"+ ask.Query, "Ans:"+ answer.Reply)]
        # chat_data.append(chat_data_ele)
        print("chat_data: ", chat_data)

if __name__=="__main__":
    ip = "127.0.0.1"
    HOST,PORT = ip,6789
    print("server@" + ip + ":start")
    TCPServer.allow_reuse_address = True
    with TCPServer((HOST,PORT), GPBHandler) as server:
        server.serve_forever()

