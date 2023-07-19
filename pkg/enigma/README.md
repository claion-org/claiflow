# enigma

## config file

```yaml
blockMethod: AES
blockSize: 128
blockKey: YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n # base64("brown fox jumps over the lazy dog")
cipherMode: GCM
cipherSalt: c2FsdHk= # base64("salty")
padding: NONE
strconv: base64
```

## usage

```go
func main() {
    NewString := func(s string) *string { return &s }
    
    config := enigma.Config{}
    config.EncryptionMethod = "AES"
    config.BlockSize = 128
    config.BlockKey = "YnJvd24gZm94IGp1bXBzIG92ZXIgdGhlIGxhenkgZG9n"
    config.CipherMode = "GCM"
    config.Padding = "PKCS"
    config.CipherSalt = NewString("64uk656M7KWQIO2XjCDss4frsJTtgLTsl5Ag7YOA6rOg7YyM")
    config.StrConv = "base64"

    const example = "세종어제 훈민정음\n" +
        "나랏말이\n" +
        "중국과 달라\n" +
        "문자와 서로 통하지 아니하므로\n" +
        "이런 까닭으로 어리석은 백성이 이르고자 하는 바가 있어도\n" +
        "마침내 제 뜻을 능히 펴지 못하는 사람이 많다.\n" +
        "내가 이를 위해 불쌍히 여겨\n" +
        "새로 스물여덟 글자를 만드니\n" +
        "사람마다 하여금 쉬이 익혀 날마다 씀에 편안케 하고자 할 따름이다.\n"

    machine, err := enigma.NewMachine(config.ToOption())
    if err != nil {
        panic(err)
    }

    encoded, err := machine.Encode([]byte(example))
    if err != nil {
        panic(err)
    }

    plain, err := machine.Decode(encoded)
    if err != nil {
        panic(err)
    }

    if !strings.EqualFold(example, string(plain)) {
        panic("diff")
    }
}
```

## table of configurations

- block

    | blockMethod | blockSize     | blockKey       |
    |---          |---            |---             |
    | NONE        | default(1)    | base64(string) |
    | *AES        | 128, 192, 256 | base64(string) |
    | DES         | 64            | base64(string) |

- cipher

    | cipherMode | blockMethod    |
    |---         |---             |
    | NONE       | NONE, AES, DES |
    | *GCM       | AES            |
    | CBC        | NONE, AES, DES |

    | cipherSalt     |
    |---             |
    | base64(string) |
    | nil            |

    > cipherSalt를 지정되지 않은 경우 암호화 하면서 생성한 salt값을 암호화 결과의 앞에 붙여서 리턴 한다.
    > 복호화에서는 앞에 붙어있는 salt값을 분리하여 사용

- padding

    | blockMethod | cipherMode | padding    |
    |---          |---         |---         |
    | AES         | NONE       | PKCS       |
    | AES         | CBC        | PKCS       |
    | AES         | GCM        | NONE, PKCS |
    | DES         | NONE       | PKCS       |
    | DES         | CBC        | PKCS       |

- strconv

    | strconv |
    |---      |
    | plain   |
    | *base64 |
    | hex     |

## 설정 순서

1. 암호화 블럭을 만든다 (blockMethod, blockSize, blockKey)

1. 암호화 cipher를 만든다 (cipherMode)

    - GCM: 나중에 nonce를 생성하기 위해서 cipher.AEAD.NonceSize() 값을 저장

    - CBC: 나중에 iv를 생성하기 위해서 cipher.Block.BlockSize() 값을 저장

1. 암호화 블럭과 cipher를 이용하여 enigma.Machine에서 이용하는 Encoder, Decoder 함수를 생성하여 enigma.Machine 생성

## 암호화 순서

```go
func (machine *Machine) Encode(src []byte) ([]byte, error)
```

1. 암화와 블럭 사이즈 만큼 입력값에 패드 추가

1. Encoder 실행

1. salt encode rule 적용; salt 값이 null이면 암호화 결과에 salt를 앞에 붙이는 작업

1. strconv encode; 지정된 변환 설정에 따라 []byte 결과를 인코드 한다

## 복호화 순서 `(조립의 역순)`

```go
func (machine *Machine) Decode(src []byte) ([]byte, error)
```

1. strconv decode; 지정된 변환 설정에 따라 []byte 결과를 디코드 한다

1. salt decode rule 적용; salt 값이 null이면 암호화 결과에서 앞에 저장된 salt를 분리하는 작업

1. Decoder 실행

1. 암화와 블럭 사이즈 만큼 입력값에 패드 제거
