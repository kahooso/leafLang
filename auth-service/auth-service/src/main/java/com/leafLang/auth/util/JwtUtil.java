package com.leafLang.auth.util;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.util.Date;

@Component
public class JwtUtil {

    @Value("${jwt.secret}")
    private String secretKey;  // Секретный ключ для подписи токена

    @Value("${jwt.expiration}")
    private Long expiration;

    public JwtUtil() { }

    public String generateToken(String username) {
        return Jwts.builder()
                .setSubject(username)
                .setIssuedAt(new Date())
                .setExpiration(new Date(System.currentTimeMillis() + expiration))
                .signWith(SignatureAlgorithm.HS512, secretKey)
                .compact();
    }

    private Claims getClaimFromToken(String token) {
        return Jwts.parser()
                .setSigningKey(secretKey)  // Устанавливаем секретный ключ для проверки подписи
                .build()
                .parseSignedClaims(token)
                .getPayload();
    }

    public String getUsernameFromToken(String token) {
        return getClaimFromToken(token).getSubject(); // get subject
    }

    public boolean isTokenExpired(String token) {
        return getClaimFromToken(token).getExpiration().before(new Date());
    }

    public boolean validateToken(String token, String username) {
        return (username.equals(getUsernameFromToken(token)) && !isTokenExpired(token));
    }
}