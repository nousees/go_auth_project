package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("3646e12ee950fce8016a4f81aa0db79addd49dd059a268813623388cbd847cfca9e4c20cd9af0e4ce08ef94575743775df581581bddaa63797c3033c5855b3afe3ccd2b937e7f3b8ab892a91191a6590d71185f8702e034b2dd2e60e74a98f1415a72cefc4a35b5a1c76564540f85413d4be8b635ad394ee7485058efcf4602c9a119dd232620d272b3d175cc70ebb74a75d0de5053f9c74e9b6e28dd2608e883d861dcce0446b52fc2c2d66c5f9617c7e1ecf8723127c7bdbe63ee391cc45cf8ac4abf40d02f4da0c0cefa0e45fa69c8c8bcf9d9f0a52975b996c1b0886700443ef444c7ff352e6d26032b201691270b6580429d464752e574522d10f98488f")

func GenerateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
