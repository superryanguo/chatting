from keras.layers import Input, Dense
from keras.models import Model

# This returns a tensor
inputs = Input(shape=(784,))

# a layer instance is callable on a tensor, and returns a tensor
x = Dense(64, activation='relu')(inputs)
# 输入inputs，输出x
# (inputs)代表输入
x = Dense(64, activation='relu')(x)
# 输入x，输出x
predictions = Dense(10, activation='softmax')(x)
# 输入x，输出分类

# This creates a model that includes
# the Input layer and three Dense layers
model = Model(inputs=inputs, outputs=predictions)
model.compile(optimizer='rmsprop',
              loss='categorical_crossentropy',
              metrics=['accuracy'])
model.fit(data, labels)  # starts training
# //可以看到结构与序贯模型完全不一样，其中x = Dense(64, activation=‘relu’)(inputs)中：(input)代表输入；x代表输出
# model = Model(inputs=inputs, outputs=predictions)；该句是函数式模型的经典，可以同时输入两个input，然后输出output两个模型
