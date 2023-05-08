// av1_wrapper.c

#include "av1_wrapper.h"
#include <libavcodec/avcodec.h>
#include <libavutil/imgutils.h>
#include <libavutil/opt.h>

int av1_encode(const uint8_t *input, size_t input_size, uint8_t **output, size_t *output_size) {
    // 这里仅作为示例，假设输入数据为YUV420格式的图像，宽度和高度分别为640和480
    int width = 640;
    int height = 480;

    AVCodec *codec;
    AVCodecContext *c = NULL;
    AVFrame *frame;
    AVPacket pkt;
    int ret;

    codec = avcodec_find_encoder_by_name("libaom-av1");
    if (!codec) {
        return -1;
    }

    c = avcodec_alloc_context3(codec);
    if (!c) {
        return -2;
    }

    c->width = width;
    c->height = height;
    c->time_base = (AVRational){1, 25};
    c->pix_fmt = AV_PIX_FMT_YUV420P;

    av_opt_set(c->priv_data, "crf", "30", 0);
    av_opt_set(c->priv_data, "strict", "experimental", 0);

    ret = avcodec_open2(c, codec, NULL);
    if (ret < 0) {
        return -3;
    }

    frame = av_frame_alloc();
    if (!frame) {
        return -4;
    }
    frame->format = c->pix_fmt;
    frame->width  = c->width;
    frame->height = c->height;

    ret = av_frame_get_buffer(frame, 32);
    if (ret < 0) {
        return -5;
    }

    // 填充YUV数据，这里仅作为示例，实际情况可能需要从input参数中获取数据
    for (int y = 0; y < height; y++) {
        for (int x = 0; x < width; x++) {
            frame->data[0][y * frame->linesize[0] + x] = x + y;
        }
    }
    for (int y = 0; y < height / 2; y++) {
        for (int x = 0; x < width / 2; x++) {
            frame->data[1][y * frame->linesize[1] + x] = 128 + y;
            frame->data[2][y * frame->linesize[2] + x] = 64 + x;
        }
    }

    av_init_packet(&pkt);
    pkt.data = NULL;
    pkt.size = 0;

    ret = avcodec_send_frame(c, frame);
    if (ret < 0) {
        return -6;
    }

    ret = avcodec_receive_packet(c, &pkt);
    if (ret < 0) {
        return -7;
    }

    *output_size = pkt.size;
    *output = (uint8_t *)malloc(pkt.size);
    if (!*output) {
        return -8;
    }
    memcpy(*output, pkt.data, pkt.size);
    av_packet_unref(&pkt);
    av_frame_free(&frame);
    avcodec_free_context(&c);

    return 0;
}

int av1_decode(const uint8_t *input, size_t input_size, uint8_t **output, size_t *output_size) {
    AVCodec *codec;
    AVCodecContext *c = NULL;
    AVFrame *frame;
    AVPacket pkt;
    int ret;
    codec = avcodec_find_decoder_by_name("libaom-av1");
    if (!codec) {
        return -1;
    }

    c = avcodec_alloc_context3(codec);
    if (!c) {
        return -2;
    }

    ret = avcodec_open2(c, codec, NULL);
    if (ret < 0) {
        return -3;
    }

    av_init_packet(&pkt);
    pkt.data = (uint8_t *)input;
    pkt.size = input_size;

    ret = avcodec_send_packet(c, &pkt);
    if (ret < 0) {
        return -4;
    }

    frame = av_frame_alloc();
    if (!frame) {
        return -5;
    }

    ret = avcodec_receive_frame(c, frame);
    if (ret < 0) {
        return -6;
    }

    int width = frame->width;
    int height = frame->height;

    *output_size = width * height * 3 / 2; // 假设输出数据为YUV420格式
    *output = (uint8_t *)malloc(*output_size);
    if (!*output) {
        return -7;
    }

    // 提取YUV数据
    for (int y = 0; y < height; y++) {
        memcpy(*output + y * width, frame->data[0] + y * frame->linesize[0], width);
    }
    for (int y = 0; y < height / 2; y++) {
        memcpy(*output + width * height + y * width / 2, frame->data[1] + y * frame->linesize[1], width / 2);
        memcpy(*output + width * height * 5 / 4 + y * width / 2, frame->data[2] + y * frame->linesize[2], width / 2);
    }

    av_frame_free(&frame);
    avcodec_free_context(&c);

    return 0;
}