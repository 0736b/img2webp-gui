// Copyright 2011 Google Inc. All Rights Reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the COPYING file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS. All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.
// -----------------------------------------------------------------------------
//
//   WebP encoder: main interface
//
// Author: Skal (pascal.massimino@gmail.com)

#ifndef WEBP_WEBP_ENCODE_H_
#define WEBP_WEBP_ENCODE_H_

#include "./types.h"

#ifdef __cplusplus
extern "C" {
#endif

#define WEBP_ENCODER_ABI_VERSION 0x020f    // MAJOR(8b) + MINOR(8b)

WEBP_EXTERN size_t WebPEncodeRGBA(const uint8_t* rgba,
                                  int width, int height, int stride,
                                  float quality_factor, uint8_t** output);

#ifdef __cplusplus
}    // extern "C"
#endif

#endif  // WEBP_WEBP_ENCODE_H_
