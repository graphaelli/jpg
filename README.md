# JPEG Structure

[So far just] an implementation of some of exiv2's printStructure in go.

## Library Installation:

```
go get -d github.com/graphaelli/jpg
```

## Command Line Usage [like exiv2 -pS]:
```
go run github.com/graphaelli/jpg file.jpg
```

outputs:
```
 address | marker     |  length | data
       2 | 0xd8 SOI   |       0 |
       4 | 0xe0 APP0  |      16 | JFIF.....H.H..
      22 | 0xe1 APP1  |    1994 | Exif..MM.*.................z.....................................
    2018 | 0xe1 APP1  |    5507 | http://ns.adobe.com/xap/1.0/.<?xpacket begin="..." id="W5M0MpCehi
    7527 | 0xed APP13 |      56 | Photoshop 3.0.8BIM........8BIM.%....................B~
    7585 | 0xe2 APP2  |     564 | ICC_PROFILE......$appl....mntrRGB XYZ ...........9acspAPPL....APP
    8151 | 0xc0 SOF0  |      17 |
    8170 | 0xc4 DHT   |      31 |
    8203 | 0xc4 DHT   |     181 |
    8386 | 0xc4 DHT   |      31 |
    8419 | 0xc4 DHT   |     181 |
    8602 | 0xdb DQT   |      67 |
    8671 | 0xdb DQT   |      67 |
    8740 | 0xdd DRI   |       4 |
    8746 | 0xda SOS   |      12 |
```
