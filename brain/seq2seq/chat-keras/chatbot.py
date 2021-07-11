from tensorflow.keras.models import load_model, Model
from tensorflow.keras.layers import Input, Concatenate
import tensorflow as tf
import os
from tensorflow.python.keras.layers import Layer
from tensorflow.python.keras import backend as K
import pickle
import numpy as np
import re
from AttentionLayer import AttentionLayer

with open('dic.pkl', 'rb') as f:
    vocab = pickle.load(f)
with open('inv.pkl', 'rb') as f:
    inv_vocab = pickle.load(f)


def clean_text(txt):
    txt = txt.lower()
    txt = re.sub(r"i'm", "i am", txt)
    txt = re.sub(r"he's", "he is", txt)
    txt = re.sub(r"she's", "she is", txt)
    txt = re.sub(r"that's", "that is", txt)
    txt = re.sub(r"what's", "what is", txt)
    txt = re.sub(r"where's", "where is", txt)
    txt = re.sub(r"\'ll", " will", txt)
    txt = re.sub(r"\'ve", " have", txt)
    txt = re.sub(r"\'re", " are", txt)
    txt = re.sub(r"\'d", " would", txt)
    txt = re.sub(r"won't", "will not", txt)
    txt = re.sub(r"can't", "can not", txt)
    txt = re.sub(r"[^\w\s]", "", txt)
    return txt



attn_layer = AttentionLayer()

model = load_model('./outModel/chatbot.h5', custom_objects={'AttentionLayer' : attn_layer})



encoder_inputs = model.layers[0].input
embed = model.layers[2]
enc_embed = embed(encoder_inputs)
enocoder_layer = model.layers[3]

encoder_outputs, fstate_h, fstate_c, bstate_h, bstate_c = enocoder_layer(enc_embed)

h = Concatenate()([fstate_h, bstate_h])
c = Concatenate()([fstate_c, bstate_c])
encoder_states = [h, c]

enc_model = Model(encoder_inputs, 
                    [encoder_outputs,
                     encoder_states])


latent_dim = 800

decoder_inputs = model.layers[1].input
decoder_lstm = model.layers[6]
decoder_dense = model.layers[9]
decoder_state_input_h = Input(shape=(latent_dim,), name='input_3')
decoder_state_input_c = Input(shape=(latent_dim,), name='input_4')

decoder_states_inputs = [decoder_state_input_h, decoder_state_input_c]

dec_embed = embed(decoder_inputs)

decoder_outputs, state_h, state_c = decoder_lstm(dec_embed, initial_state=decoder_states_inputs)
decoder_states = [state_h, state_c]

dec_model = Model([decoder_inputs, decoder_states_inputs], [decoder_outputs] + decoder_states)

dec_dense = model.layers[-1]
attn_layer = model.layers[7]

from tensorflow.keras.preprocessing.sequence import pad_sequences

def ai_response(query):
        try:
            query = clean_text(query)
            prepro = [query]
 
            txt = []
            for x in prepro:
                lst = []
                for y in x.split():
                    try:
                        lst.append(vocab[y])
                    except:
                        lst.append(vocab['<OUT>'])
                txt.append(lst)
            txt = pad_sequences(txt, 13, padding='post')


            ###
            enc_op, stat = enc_model.predict( txt )

            empty_target_seq = np.zeros( ( 1 , 1) )
            empty_target_seq[0, 0] = vocab['<SOS>']
            stop_condition = False
            decoded_translation = ''


            while not stop_condition :

                dec_outputs , h , c = dec_model.predict([ empty_target_seq ] + stat )

                ###
                ###########################
                attn_op, attn_state = attn_layer([enc_op, dec_outputs])
                decoder_concat_input = Concatenate(axis=-1)([dec_outputs, attn_op])
                decoder_concat_input = dec_dense(decoder_concat_input)
                ###########################

                sampled_word_index = np.argmax( decoder_concat_input[0, -1, :] )

                sampled_word = inv_vocab[sampled_word_index] + ' '

                if sampled_word != '<EOS> ':
                    decoded_translation += sampled_word


                if sampled_word == '<EOS> ' or len(decoded_translation.split()) > 13:
                    stop_condition = True

                empty_target_seq = np.zeros( ( 1 , 1 ) )
                empty_target_seq[ 0 , 0 ] = sampled_word_index
                stat = [ h , c ]

            reponse="ChatBot: "+ decoded_translation
            # print("chatbot attention : ", decoded_translation )
            # print("==============================================")

        except:
            reponse="sorry didn't got you , please type again :( "
            # print("sorry didn't got you , please type again :( ")

        return reponse

# print("##########################################")
# print("#       start chatting ver. 1.0          #")
# print("##########################################")
# prepro1 = ""
# while prepro1 != 'q':
        # prepro1 = input("you : ")
        # print(ai_response(prepro1))

#add the server handling for GPB message
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

        answer.Reply =  ai_response(ask.Query)

        # if ask.Query == "How are you?" :
            # answer.Reply =  "I'm Fine. Hope you have a good day!"
        # else :
            # answer.Reply =  "Hello."

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
    ip = "0.0.0.0"
    HOST,PORT = ip,8099
    print("server@" + ip + ":start")
    TCPServer.allow_reuse_address = True
    with TCPServer((HOST,PORT), GPBHandler) as server:
        server.serve_forever()

