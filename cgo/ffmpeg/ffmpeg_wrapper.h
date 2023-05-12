#ifndef FFMPEG_WRAPPER_H
#define FFMPEG_WRAPPER_H

#include <stdint.h>

int encode_av1(uint8_t *input, int input_size, uint8_t **output, int *output_size);
int decode_av1(uint8_t *input, int input_size, uint8_t **output, int *output_size);

#endif // FFMPEG_WRAPPER_H
