#include <libavcodec/avcodec.h>
#include <libavutil/opt.h>
#include <stdlib.h>
#include "ffmpeg_wrapper.h"

// AV1编码函数
int encode_av1(uint8_t *input, int input_size, uint8_t **output, int *output_size) {
    AVCodec *codec;
    AVCodecContext *c = NULL;
    int ret, i;
    AVFrame *frame;
    AVPacket *pkt;

    // 查找AV1编码器
    codec = avcodec_find_encoder(AV_CODEC_ID_AV1);
    if (!codec) {
        return -1; // 无法找到编码器
    }

    c = avcodec_alloc_context3(codec);
    if (!c) {
        return -2; // 无法分配编码器上下文
    }

    // 设置编码器参数
    c->width = 1920;
    c->height = 1080;
    c->time_base = (AVRational){1, 25};
    c->pix_fmt = AV_PIX_FMT_YUV420P;

    // 打开编码器
    if (avcodec_open2(c, codec, NULL) < 0) {
        return -3; // 无法打开编码器
    }

    frame = av_frame_alloc();
    if (!frame) {
        return -4; // 无法分配帧
    }

    frame->format = c->pix_fmt;
    frame->width  = c->width;
    frame->height = c->height;

    ret = av_frame_get_buffer(frame, 0);
    if (ret < 0) {
        return -5; // 无法分配帧数据
    }

    pkt = av_packet_alloc();
    if (!pkt) {
        return -6; // 无法分配数据包
    }

    // 在这里，我们假设输入数据是YUV420P格式的，适合我们的帧
    // 注意:在实际应用中，你可能需要根据你的需求来处理输入数据
    memcpy(frame->data[0], input, input_size);

    // 发送帧到编码器
    ret = avcodec_send_frame(c, frame);
    if (ret < 0) {
        return -7; // 发送帧失败
    }

    // 获取编码的数据包
    ret = avcodec_receive_packet(c, pkt);
    if (ret < 0) {
        return -8; // 获取数据包失败
    }

    // 将编码的数据包复制到输出
    *output = (uint8_t *)malloc(pkt->size);
    if (!*output) {
        return -9; // 分配输出内存失败
    }

    memcpy(*output, pkt->data, pkt->size);
    *output_size = pkt->size;

    // 清理
    av_packet_free(&pkt);
    av_frame_free(&frame);
    avcodec_free_context(&c);

    return 0; // 成功
}

// AV1解码函数
int decode_av1(uint8_t *input, int input_size, uint8_t **output, int *output_size) {
    AVCodec *codec;
    AVCodecContext *c = NULL;
    int ret;
    AVFrame *frame;
    AVPacket *pkt;

    // 查找AV1解码器
    codec = avcodec_find_decoder(AV_CODEC_ID_AV1);
    if (!codec) {
        return -1; // 无法找到解码器
    }

    c = avcodec_alloc_context3(codec);
    if (!c) {
        return -2; // 无法分配解码器上下文
    }

    // 打开解码器
    if (avcodec_open2(c, codec, NULL) < 0) {
        return -3; // 无法打开解码器
    }

    frame = av_frame_alloc();
    if (!frame) {
        return -4; // 无法分配帧
    }

    pkt = av_packet_alloc();
    if (!pkt) {
        return -5; // 无法分配数据包
    }

    // 在这里，我们假设输入数据是AV1编码的
    // 注意:在实际应用中，你可能需要根据你的需求来处理输入数据
    pkt->data = input;
    pkt->size = input_size;

    // 发送数据包到解码器
    ret = avcodec_send_packet(c, pkt);
    if (ret < 0) {
        return -6; // 发送数据包失败
    }

    // 获取解码的帧
    ret = avcodec_receive_frame(c, frame);
    if (ret < 0) {
        return -7; // 获取帧失败
    }

    // 将解码的帧复制到输出
    // 注意:这只是一个简单的例子，我们只处理了YUV420P格式的帧
    // 在实际应用中，你可能需要根据你的需求来处理解码的帧
    *output_size = frame->linesize[0] * c->height;
    *output = (uint8_t *)malloc(*output_size);
    if (!*output) {
        return -8; // 分配输出内存失败
    }

    memcpy(*output, frame->data[0], *output_size);

    // 清理
    av_packet_free(&pkt);
    av_frame_free(&frame);
    avcodec_free_context(&c);

    return 0; // 成功
}

