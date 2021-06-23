from keras.layers import Input, Embedding, LSTM, Dense
from keras.models import Model

# Headline input: meant to receive sequences of 100 integers, between 1 and 10000.
# Note that we can name any layer by passing it a "name" argument.
main_input = Input(shape=(100,), dtype='int32', name='main_input')
# 一个100词的BOW序列

# This embedding layer will encode the input sequence
# into a sequence of dense 512-dimensional vectors.
x = Embedding(output_dim=512, input_dim=10000, input_length=100)(main_input)
# Embedding层，把100维度再encode成512的句向量，10000指的是词典单词总数


# A LSTM will transform the vector sequence into a single vector,
# containing information about the entire sequence
lstm_out = LSTM(32)(x)
# ？ 32什么意思？？？？？？？？？？？？？？？？？？？？？

#然后，我们插入一个额外的损失，使得即使在主损失很高的情况下，LSTM和Embedding层也可以平滑的训练。

auxiliary_output = Dense(1, activation='sigmoid', name='aux_output')(lstm_out)
#再然后，我们将LSTM与额外的输入数据串联起来组成输入，送入模型中：
# 模型一：只针对以上的序列做的预测模型
# 模型二：组合模型
auxiliary_input = Input(shape=(5,), name='aux_input')  # 新加入的一个Input,5维度
x = keras.layers.concatenate([lstm_out, auxiliary_input])   # 组合起来，对应起来


# We stack a deep densely-connected network on top
# 组合模型的形式
x = Dense(64, activation='relu')(x)
x = Dense(64, activation='relu')(x)
x = Dense(64, activation='relu')(x)
# And finally we add the main logistic regression layer
main_output = Dense(1, activation='sigmoid', name='main_output')(x)


#最后，我们定义整个2输入，2输出的模型：
model = Model(inputs=[main_input, auxiliary_input], outputs=[main_output, auxiliary_output])
#模型定义完毕，下一步编译模型。
#我们给额外的损失赋0.2的权重。我们可以通过关键字参数loss_weights或loss来为不同的输出设置不同的损失函数或权值。
#这两个参数均可为Python的列表或字典。这里我们给loss传递单个损失函数，这个损失函数会被应用于所有输出上。

# 训练方式一：两个模型一个loss
model.compile(optimizer='rmsprop', loss='binary_crossentropy',
              loss_weights=[1., 0.2])
#编译完成后，我们通过传递训练数据和目标值训练该模型：

model.fit([headline_data, additional_data], [labels, labels],
          epochs=50, batch_size=32)

# 训练方式二：两个模型,两个Loss
#因为我们输入和输出是被命名过的（在定义时传递了“name”参数），我们也可以用下面的方式编译和训练模型：
model.compile(optimizer='rmsprop',
              loss={'main_output': 'binary_crossentropy', 'aux_output': 'binary_crossentropy'},
              loss_weights={'main_output': 1., 'aux_output': 0.2})

# And trained it via:
model.fit({'main_input': headline_data, 'aux_input': additional_data},
          {'main_output': labels, 'aux_output': labels},
          epochs=50, batch_size=32)

