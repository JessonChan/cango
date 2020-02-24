// Copyright 2020 Cango Author.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var globalFiles = map[string][]byte{}
var cmdArr = []string{"demo"}

func init() {
	globalFiles["./demo"] = []byte{}
	globalFiles["./demo/controller"] = []byte{}
	globalFiles["./demo/controller/can.go"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 132, 145, 65, 75, 244, 48, 16, 134, 207, 157, 95, 49, 95, 14, 31, 173, 44, 237, 189, 139, 7, 217, 139, 43, 120, 89, 92, 47, 34, 110, 204, 206, 182, 209, 52, 9, 201, 116, 85, 74, 255, 187, 164, 237, 65, 16, 214, 75, 2, 201, 51, 239, 251, 132, 120, 169, 222, 101, 67, 168, 156, 229, 224, 140, 161, 0, 160, 59, 239, 2, 99, 14, 153, 104, 52, 183, 253, 107, 169, 92, 87, 221, 81, 140, 206, 110, 90, 105, 43, 37, 109, 227, 4, 20, 0, 252, 229, 9, 55, 210, 110, 56, 24, 140, 28, 122, 197, 56, 64, 54, 1, 229, 126, 183, 197, 195, 89, 154, 158, 106, 81, 137, 3, 140, 0, 103, 25, 240, 5, 175, 113, 6, 118, 212, 232, 200, 20, 246, 187, 109, 254, 127, 73, 25, 198, 2, 224, 212, 91, 133, 185, 194, 171, 229, 176, 192, 173, 61, 210, 103, 238, 227, 229, 146, 117, 165, 19, 55, 175, 101, 203, 157, 73, 181, 5, 106, 203, 20, 78, 82, 209, 48, 166, 201, 64, 220, 7, 187, 72, 220, 187, 35, 153, 71, 77, 31, 3, 100, 217, 131, 55, 53, 34, 138, 31, 1, 43, 200, 178, 137, 169, 81, 220, 146, 49, 110, 53, 205, 253, 75, 23, 99, 122, 211, 111, 217, 27, 175, 255, 84, 149, 94, 175, 165, 215, 229, 91, 116, 246, 146, 100, 39, 253, 83, 228, 160, 109, 243, 60, 111, 73, 83, 180, 73, 68, 212, 40, 230, 175, 88, 84, 190, 3, 0, 0, 255, 255, 203, 91, 31, 210, 205, 1, 0, 0}
	globalFiles["./demo/filter"] = []byte{}
	globalFiles["./demo/filter/visit.go"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 116, 143, 193, 106, 243, 48, 16, 132, 207, 222, 167, 88, 124, 248, 177, 195, 143, 252, 4, 61, 185, 148, 82, 90, 8, 134, 246, 90, 20, 101, 45, 139, 42, 146, 179, 90, 5, 138, 241, 187, 23, 217, 61, 228, 210, 155, 180, 59, 223, 206, 204, 172, 205, 151, 182, 132, 163, 243, 66, 12, 224, 46, 115, 100, 193, 6, 170, 58, 144, 116, 147, 200, 92, 3, 84, 93, 135, 181, 117, 50, 229, 147, 50, 241, 210, 189, 80, 74, 49, 244, 147, 14, 157, 209, 193, 71, 91, 52, 127, 11, 108, 172, 161, 5, 144, 239, 153, 240, 195, 37, 39, 79, 155, 27, 38, 225, 108, 4, 23, 168, 54, 145, 218, 199, 155, 155, 137, 65, 56, 122, 255, 251, 61, 244, 58, 244, 194, 30, 86, 128, 155, 102, 252, 196, 7, 220, 153, 129, 172, 75, 66, 188, 179, 205, 191, 187, 251, 203, 218, 2, 140, 57, 24, 108, 110, 120, 184, 91, 180, 120, 100, 122, 214, 225, 236, 169, 97, 186, 226, 161, 212, 84, 3, 93, 51, 37, 105, 209, 5, 33, 30, 181, 161, 101, 45, 217, 74, 154, 173, 164, 234, 117, 120, 164, 83, 182, 5, 82, 111, 36, 83, 60, 255, 199, 242, 126, 31, 94, 213, 81, 203, 212, 66, 197, 36, 153, 3, 10, 103, 130, 21, 126, 2, 0, 0, 255, 255, 23, 72, 165, 200, 95, 1, 0, 0}
	globalFiles["./demo/go.mod"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 124, 204, 49, 14, 194, 48, 12, 133, 225, 61, 167, 200, 9, 226, 148, 29, 22, 54, 118, 14, 144, 186, 150, 19, 148, 198, 197, 78, 144, 184, 61, 162, 123, 59, 255, 223, 123, 171, 44, 163, 146, 231, 210, 243, 152, 3, 202, 10, 15, 50, 147, 118, 207, 169, 193, 252, 237, 100, 42, 181, 150, 198, 206, 177, 248, 41, 76, 23, 231, 148, 222, 163, 232, 209, 8, 83, 99, 241, 159, 24, 98, 136, 127, 187, 213, 132, 231, 246, 122, 243, 240, 52, 82, 131, 215, 158, 112, 79, 178, 16, 176, 192, 150, 122, 6, 83, 132, 147, 7, 247, 11, 0, 0, 255, 255, 35, 59, 45, 223, 198, 0, 0, 0}
	globalFiles["./demo/go.sum"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 164, 207, 79, 111, 178, 48, 28, 0, 224, 187, 159, 194, 187, 17, 126, 237, 43, 20, 223, 196, 131, 248, 23, 17, 137, 64, 88, 217, 173, 45, 72, 49, 163, 149, 32, 176, 249, 233, 151, 237, 180, 29, 118, 242, 19, 60, 121, 202, 234, 46, 59, 110, 8, 93, 155, 135, 162, 109, 181, 90, 73, 166, 76, 193, 212, 155, 46, 199, 61, 24, 96, 192, 20, 3, 6, 192, 24, 193, 220, 178, 48, 158, 254, 67, 136, 16, 142, 68, 97, 11, 24, 75, 244, 223, 63, 205, 37, 137, 9, 77, 242, 151, 124, 29, 183, 30, 100, 183, 46, 57, 139, 166, 91, 110, 136, 187, 161, 22, 17, 65, 4, 101, 134, 86, 82, 44, 70, 79, 122, 102, 169, 141, 90, 231, 95, 236, 190, 137, 214, 71, 241, 177, 241, 211, 135, 231, 62, 84, 19, 135, 245, 225, 170, 216, 118, 82, 238, 219, 112, 235, 250, 75, 82, 177, 196, 91, 207, 46, 116, 249, 23, 123, 109, 59, 245, 27, 5, 27, 102, 200, 1, 152, 2, 118, 138, 11, 114, 44, 100, 243, 111, 77, 196, 236, 28, 136, 222, 230, 89, 46, 82, 185, 117, 233, 177, 210, 50, 229, 93, 29, 36, 133, 162, 175, 3, 61, 134, 39, 110, 178, 247, 161, 210, 79, 106, 63, 138, 161, 143, 205, 230, 182, 195, 121, 42, 206, 247, 213, 224, 98, 159, 202, 97, 23, 112, 238, 245, 19, 39, 137, 28, 75, 161, 147, 223, 23, 145, 156, 103, 139, 209, 103, 0, 0, 0, 255, 255, 100, 103, 207, 171, 202, 1, 0, 0}
	globalFiles["./demo/main.go"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 164, 83, 77, 111, 227, 54, 16, 61, 147, 191, 98, 150, 135, 133, 20, 184, 210, 30, 23, 46, 122, 8, 92, 164, 235, 98, 55, 107, 216, 73, 138, 34, 8, 26, 70, 30, 203, 108, 40, 146, 161, 70, 249, 18, 244, 223, 139, 145, 100, 199, 65, 154, 246, 176, 23, 129, 122, 124, 51, 111, 230, 113, 38, 232, 226, 86, 151, 8, 149, 54, 78, 74, 83, 5, 31, 9, 18, 41, 148, 67, 202, 183, 68, 65, 73, 161, 124, 173, 164, 20, 121, 14, 170, 52, 180, 109, 110, 178, 194, 87, 249, 239, 88, 215, 222, 205, 182, 218, 229, 133, 118, 214, 151, 204, 249, 1, 66, 233, 149, 76, 165, 164, 167, 128, 48, 211, 110, 70, 209, 66, 77, 177, 41, 8, 90, 41, 122, 66, 118, 190, 156, 195, 245, 189, 182, 13, 78, 85, 174, 174, 101, 39, 229, 189, 142, 240, 23, 252, 2, 3, 97, 137, 165, 169, 9, 227, 249, 114, 158, 124, 28, 179, 180, 93, 42, 229, 166, 113, 5, 36, 5, 28, 141, 96, 10, 115, 183, 198, 199, 36, 212, 255, 45, 242, 115, 110, 152, 55, 124, 179, 45, 85, 150, 101, 83, 48, 142, 48, 110, 116, 129, 109, 199, 145, 17, 169, 137, 110, 44, 226, 155, 95, 163, 189, 48, 248, 208, 74, 33, 206, 130, 157, 2, 128, 58, 72, 48, 145, 66, 244, 156, 41, 168, 47, 104, 173, 159, 244, 113, 31, 248, 162, 227, 158, 122, 15, 22, 252, 44, 47, 181, 125, 223, 108, 106, 36, 214, 149, 98, 101, 158, 17, 160, 63, 119, 255, 214, 218, 113, 48, 255, 219, 88, 235, 116, 133, 93, 222, 62, 161, 142, 93, 174, 131, 201, 254, 174, 189, 83, 215, 59, 242, 194, 215, 244, 13, 105, 235, 215, 82, 156, 234, 10, 129, 211, 25, 87, 74, 241, 39, 234, 56, 136, 139, 153, 183, 62, 238, 47, 184, 226, 247, 189, 169, 116, 184, 28, 152, 87, 7, 4, 182, 72, 113, 41, 106, 10, 16, 234, 140, 165, 216, 31, 165, 203, 30, 98, 140, 5, 123, 172, 96, 57, 53, 101, 172, 87, 238, 193, 218, 60, 239, 130, 217, 151, 87, 30, 94, 152, 218, 208, 137, 177, 132, 241, 141, 27, 3, 220, 143, 117, 225, 29, 69, 111, 45, 255, 238, 92, 124, 127, 182, 134, 192, 228, 227, 65, 242, 131, 17, 187, 135, 163, 131, 139, 20, 22, 17, 191, 104, 183, 182, 152, 68, 188, 131, 35, 94, 169, 108, 137, 119, 13, 214, 244, 198, 169, 97, 77, 178, 153, 118, 191, 226, 77, 83, 114, 68, 54, 188, 193, 4, 248, 124, 190, 252, 154, 45, 52, 109, 211, 189, 169, 20, 27, 220, 207, 0, 239, 112, 146, 114, 34, 238, 105, 200, 181, 66, 250, 35, 26, 174, 119, 4, 78, 241, 225, 196, 88, 28, 65, 149, 83, 21, 134, 237, 251, 73, 135, 144, 241, 150, 166, 19, 80, 199, 33, 168, 116, 95, 208, 75, 18, 95, 103, 43, 90, 251, 134, 94, 113, 74, 159, 205, 157, 161, 175, 190, 44, 95, 132, 126, 195, 29, 144, 102, 99, 116, 186, 167, 75, 33, 78, 241, 97, 166, 93, 146, 242, 121, 137, 229, 89, 176, 39, 141, 43, 18, 117, 211, 16, 121, 119, 193, 99, 170, 38, 192, 141, 37, 233, 56, 98, 220, 154, 216, 181, 174, 10, 107, 138, 91, 168, 240, 131, 146, 66, 116, 67, 158, 198, 37, 131, 192, 241, 122, 29, 219, 133, 143, 52, 133, 207, 159, 62, 125, 238, 38, 227, 35, 126, 15, 84, 183, 103, 193, 174, 154, 205, 198, 60, 78, 225, 242, 106, 72, 221, 170, 97, 55, 187, 9, 244, 222, 247, 107, 203, 238, 118, 169, 236, 228, 63, 1, 0, 0, 255, 255, 76, 79, 200, 11, 38, 5, 0, 0}
	globalFiles["./demo/manager"] = []byte{}
	globalFiles["./demo/model"] = []byte{}
	globalFiles["./demo/static"] = []byte{}
	globalFiles["./demo/static/css"] = []byte{}
	globalFiles["./demo/static/css/index.css"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 210, 75, 73, 44, 202, 78, 202, 41, 77, 85, 168, 230, 82, 80, 80, 80, 72, 206, 207, 201, 47, 178, 82, 128, 137, 90, 115, 213, 114, 113, 37, 90, 229, 100, 230, 101, 43, 84, 67, 216, 25, 249, 101, 169, 69, 168, 170, 139, 82, 147, 82, 147, 147, 19, 11, 74, 139, 10, 114, 64, 90, 0, 1, 0, 0, 255, 255, 254, 195, 9, 179, 85, 0, 0, 0}
	globalFiles["./demo/static/favicon.ico"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 0, 117, 2, 138, 253, 137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 16, 0, 0, 0, 16, 8, 6, 0, 0, 0, 31, 243, 255, 97, 0, 0, 2, 60, 73, 68, 65, 84, 56, 141, 133, 147, 91, 72, 147, 97, 0, 134, 63, 195, 234, 34, 130, 168, 46, 130, 48, 130, 192, 240, 38, 161, 155, 174, 138, 32, 232, 50, 10, 247, 167, 81, 33, 22, 166, 86, 12, 65, 49, 148, 148, 130, 186, 200, 80, 8, 67, 178, 130, 200, 74, 230, 191, 50, 115, 182, 245, 87, 104, 82, 90, 233, 194, 169, 149, 121, 152, 139, 153, 90, 185, 105, 30, 182, 166, 255, 225, 233, 98, 120, 154, 154, 31, 60, 23, 223, 233, 249, 224, 227, 125, 69, 130, 44, 10, 76, 178, 8, 73, 178, 96, 134, 147, 85, 235, 57, 42, 175, 98, 254, 90, 36, 38, 89, 132, 76, 149, 34, 95, 72, 178, 152, 114, 116, 151, 96, 182, 239, 228, 81, 123, 46, 254, 224, 0, 0, 193, 233, 113, 110, 126, 76, 94, 81, 34, 36, 89, 112, 171, 37, 21, 221, 208, 0, 104, 240, 148, 115, 189, 241, 8, 63, 198, 58, 233, 245, 183, 252, 87, 32, 201, 2, 33, 201, 130, 103, 157, 133, 0, 168, 250, 52, 158, 209, 86, 92, 67, 10, 33, 53, 64, 69, 123, 222, 202, 130, 116, 91, 12, 170, 62, 197, 120, 200, 199, 169, 234, 205, 228, 215, 237, 165, 168, 201, 68, 150, 178, 107, 222, 193, 40, 178, 149, 120, 114, 94, 238, 38, 211, 17, 71, 113, 83, 226, 156, 160, 162, 61, 15, 0, 165, 167, 116, 201, 23, 204, 246, 88, 250, 70, 90, 1, 240, 5, 250, 241, 5, 188, 232, 134, 198, 249, 231, 59, 194, 130, 214, 161, 23, 0, 75, 126, 88, 242, 211, 13, 12, 7, 188, 24, 134, 129, 165, 163, 128, 36, 107, 52, 253, 99, 95, 0, 48, 219, 99, 195, 2, 239, 159, 207, 0, 20, 190, 59, 188, 72, 240, 176, 237, 2, 0, 175, 220, 183, 145, 100, 65, 90, 205, 86, 12, 195, 96, 44, 52, 140, 36, 71, 133, 5, 110, 191, 19, 128, 114, 87, 214, 130, 203, 73, 214, 53, 204, 236, 93, 105, 56, 136, 36, 11, 110, 124, 56, 14, 64, 163, 183, 146, 99, 143, 215, 98, 235, 42, 70, 84, 125, 189, 10, 192, 72, 112, 144, 203, 111, 14, 144, 110, 139, 225, 218, 219, 67, 244, 250, 155, 153, 25, 101, 45, 103, 200, 116, 196, 49, 56, 222, 13, 64, 105, 243, 105, 234, 251, 238, 161, 244, 148, 34, 78, 60, 89, 135, 107, 72, 33, 114, 124, 31, 109, 163, 230, 91, 209, 236, 220, 48, 12, 166, 180, 32, 0, 206, 1, 27, 63, 39, 122, 73, 169, 222, 24, 206, 129, 36, 11, 114, 95, 239, 161, 204, 153, 198, 29, 103, 6, 5, 117, 251, 102, 163, 124, 169, 126, 63, 15, 92, 57, 220, 119, 101, 241, 107, 194, 13, 128, 219, 239, 36, 163, 118, 219, 92, 144, 150, 35, 91, 137, 167, 203, 247, 30, 221, 208, 208, 116, 21, 0, 71, 79, 9, 137, 214, 232, 133, 73, 92, 14, 75, 199, 69, 126, 79, 122, 48, 219, 99, 105, 244, 90, 208, 116, 149, 179, 181, 219, 23, 38, 49, 178, 137, 243, 73, 173, 217, 130, 219, 239, 68, 211, 85, 254, 170, 147, 220, 253, 116, 110, 113, 153, 18, 44, 139, 235, 28, 73, 74, 245, 38, 146, 172, 171, 151, 172, 243, 63, 171, 252, 163, 137, 13, 228, 3, 54, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130, 1, 0, 0, 255, 255, 109, 58, 46, 156, 117, 2, 0, 0}
	globalFiles["./demo/static/js"] = []byte{}
	globalFiles["./demo/static/js/index.js"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 108, 144, 65, 79, 2, 49, 16, 133, 239, 251, 43, 38, 61, 64, 55, 97, 145, 24, 79, 154, 61, 168, 49, 65, 131, 98, 132, 131, 215, 178, 157, 149, 134, 50, 173, 237, 212, 101, 99, 248, 239, 134, 69, 20, 13, 61, 77, 250, 190, 190, 247, 58, 117, 162, 138, 141, 35, 136, 75, 215, 76, 12, 173, 100, 14, 159, 25, 0, 128, 118, 85, 90, 35, 241, 240, 13, 249, 206, 226, 110, 188, 105, 239, 181, 20, 214, 208, 74, 228, 195, 136, 124, 205, 28, 204, 34, 49, 74, 17, 185, 181, 40, 6, 32, 180, 137, 222, 170, 246, 114, 97, 93, 181, 18, 121, 182, 205, 178, 159, 140, 128, 239, 9, 35, 63, 204, 166, 79, 50, 5, 59, 0, 173, 88, 29, 242, 62, 84, 128, 37, 179, 127, 217, 67, 80, 2, 97, 3, 175, 143, 147, 241, 239, 165, 204, 175, 58, 246, 136, 27, 58, 143, 36, 251, 207, 211, 217, 188, 63, 128, 206, 181, 86, 54, 226, 9, 50, 34, 127, 143, 99, 84, 26, 131, 20, 183, 142, 24, 137, 11, 110, 125, 87, 94, 121, 111, 77, 165, 118, 101, 207, 54, 69, 211, 52, 69, 237, 194, 186, 72, 193, 34, 85, 78, 163, 22, 39, 109, 73, 203, 238, 39, 123, 205, 212, 32, 143, 245, 128, 74, 183, 51, 86, 140, 80, 150, 37, 92, 64, 175, 247, 247, 61, 43, 78, 177, 211, 206, 71, 163, 195, 58, 118, 71, 89, 12, 252, 207, 43, 122, 71, 17, 231, 184, 225, 188, 227, 182, 217, 246, 43, 0, 0, 255, 255, 140, 216, 255, 5, 196, 1, 0, 0}
	globalFiles["./demo/util"] = []byte{}
	globalFiles["./demo/view"] = []byte{}
	globalFiles["./demo/view/index.html"] = []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 108, 82, 77, 111, 19, 49, 16, 189, 231, 87, 184, 35, 84, 165, 18, 196, 105, 79, 40, 90, 175, 132, 10, 72, 84, 136, 84, 80, 144, 56, 161, 137, 119, 146, 117, 226, 216, 139, 103, 150, 52, 68, 251, 223, 145, 227, 180, 57, 208, 147, 53, 59, 111, 159, 223, 135, 171, 139, 247, 243, 219, 135, 159, 247, 31, 84, 43, 91, 95, 143, 170, 124, 40, 143, 97, 101, 128, 2, 228, 15, 132, 77, 61, 82, 74, 169, 106, 75, 130, 202, 182, 152, 152, 196, 192, 247, 135, 143, 111, 222, 194, 105, 37, 78, 60, 213, 22, 195, 42, 86, 186, 12, 101, 225, 93, 216, 168, 68, 222, 0, 203, 222, 19, 183, 68, 2, 170, 77, 180, 52, 160, 89, 80, 156, 213, 150, 89, 187, 208, 208, 227, 196, 50, 63, 49, 178, 77, 174, 19, 197, 201, 158, 129, 235, 39, 220, 154, 65, 201, 190, 35, 3, 66, 143, 162, 215, 248, 7, 11, 30, 206, 250, 122, 89, 102, 125, 149, 46, 155, 122, 84, 233, 226, 165, 90, 196, 102, 159, 143, 94, 36, 6, 21, 131, 245, 206, 110, 12, 112, 27, 119, 159, 93, 216, 140, 175, 160, 62, 28, 202, 246, 7, 250, 158, 134, 161, 210, 101, 172, 71, 135, 195, 171, 62, 249, 153, 1, 125, 52, 171, 111, 166, 55, 83, 141, 157, 155, 172, 57, 6, 24, 134, 12, 104, 80, 112, 102, 192, 70, 31, 147, 73, 212, 92, 198, 229, 50, 75, 186, 190, 100, 247, 151, 204, 245, 52, 227, 254, 187, 62, 209, 239, 158, 88, 238, 190, 205, 191, 140, 203, 45, 195, 240, 250, 196, 54, 12, 87, 207, 185, 116, 24, 148, 245, 200, 156, 127, 89, 144, 181, 216, 245, 169, 243, 4, 245, 252, 171, 58, 114, 169, 45, 93, 84, 58, 3, 107, 245, 238, 254, 147, 169, 199, 89, 242, 51, 231, 81, 222, 153, 120, 116, 54, 87, 161, 114, 141, 129, 220, 25, 40, 193, 180, 202, 57, 254, 90, 120, 204, 115, 233, 172, 21, 233, 102, 90, 239, 118, 187, 201, 202, 73, 219, 47, 38, 54, 110, 245, 29, 49, 199, 112, 219, 98, 40, 177, 192, 203, 10, 179, 131, 227, 51, 48, 208, 56, 238, 60, 238, 103, 42, 196, 64, 57, 239, 73, 78, 25, 115, 75, 167, 122, 244, 241, 69, 254, 11, 0, 0, 255, 255, 137, 159, 99, 141, 161, 2, 0, 0}
} // end

func main() {
	if len(os.Args) == 1 {
		showHelp()
		return
	}
	cmdName := strings.ToLower(os.Args[1])
	if exeCmd(cmdName) {
		return
	}
	if cmdName == "bootstrap" {
		bootstrap()
		return
	}
	showHelp()
}

func showHelp() {
	fmt.Println("cango-cli demo")
}

func bootstrap() {
	files := map[string]string{}
	var fileNames []string
	for _, cmdName := range cmdArr {
		_ = filepath.Walk("./"+cmdName, func(path string, info os.FileInfo, err error) error {
			name := splash(path) + path
			fileNames = append(fileNames, name)
			files[name] = bytesFormat(func() []byte {
				if info.IsDir() {
					return []byte{}
				}
				bd := strings.Builder{}
				wr, err := gzip.NewWriterLevel(&bd, gzip.BestCompression)
				if err != nil {
					log.Println(err)
				}
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
				}
				_, err = wr.Write(bs)
				_ = wr.Close()
				if err != nil {
					log.Println(err)
				}
				return []byte(bd.String())
			}())
			return nil
		})
		sort.Strings(fileNames)
	}
	selfName := "./cango-cli.go"
	self := func() string {
		self, err := ioutil.ReadFile(selfName)
		if err != nil {
			log.Println(err)
			return ""
		}
		return string(self)
	}()
	if self == "" {
		return
	}
	prefixMarker := "func init() {"
	suffixMarker := "} // end"
	spaceCount := 1
	globalFilesName := "globalFiles"
	selfPrefix := strings.Index(self, prefixMarker)
	selfSuffix := strings.Index(self, suffixMarker)
	content := self[:selfPrefix] + func() string {
		builder := prefixMarker + "\n"
		for _, name := range fileNames {
			builder = builder + repeat(spaceCount) + globalFilesName + "[\"" + name + "\"] = " + files[name] + "\n"
		}
		return builder
	}() + self[selfSuffix:]
	err := ioutil.WriteFile(selfName, []byte(content), 0666)
	if err != nil {
		log.Println(err)
	}
}
func exeCmd(cmd string) (matched bool) {
	var names []string
	for k, _ := range globalFiles {
		if strings.Contains(k, "./"+cmd) {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return false
	}
	sort.Strings(names)
	for _, name := range names {
		path := name
		dir := "./" + filepath.Dir(path)
		base := filepath.Base(path)
		isDir := !strings.Contains(base, ".")
		if isDir {
			dir = path
		}
		fileInfo, _ := os.Stat(dir)
		if fileInfo == nil {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("create  dir", dir)
			}
		}
		if isDir {
			continue
		}
		newFile, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		log.Println("create file", path)
		_, err = newFile.Write(func() []byte {
			bs, _ := ioutil.ReadAll(func() io.Reader {
				rd, _ := gzip.NewReader(bytes.NewReader(globalFiles[name]))
				return rd
			}())
			return bs
		}())
		if err != nil {
			log.Println(err)
		}
		_ = newFile.Close()
	}
	return true
}
func bytesFormat(bs []byte) string {
	if len(bs) == 0 {
		return "[]byte{}"
	}
	var builder = "[]byte{"
	for _, b := range bs {
		builder = builder + (fmt.Sprintf("%d, ", b))
	}
	builder = builder[0:len(builder)-2] + "}"
	return builder
}
func repeat(c int) (rs string) {
	for i := 0; i < c; i++ {
		rs = rs + "\t"
	}
	return rs
}

func splash(path string) string {
	if strings.Contains(path, "./") {
		return ""
	}
	return "./"
}
