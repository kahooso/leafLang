package com.leafLang.auth.service;

import com.leafLang.auth.dto.LoginRequest;
import com.leafLang.auth.dto.RegisterRequest;
import com.leafLang.auth.entity.User;
import com.leafLang.auth.repository.UserRepository;
import com.leafLang.auth.util.JwtUtil;
import lombok.extern.log4j.Log4j;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.Optional;

import static java.rmi.server.LogStream.log;

@Slf4j
@Service
public class AuthService {

    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtUtil jwtUtil;

    @Autowired
    public AuthService(UserRepository userRepository, PasswordEncoder passwordEncoder, JwtUtil jwtUtil) {
        this.userRepository = userRepository;
        this.passwordEncoder = passwordEncoder;
        this.jwtUtil = jwtUtil;
    }

    public String register(RegisterRequest request) {
        if (userRepository.findByUsername(request.getUsername()) != null) {
            throw new IllegalArgumentException("Username is already in use");
        }

        User user = new User();
        user.setUsername(request.getUsername());
        user.setPassword(passwordEncoder.encode(request.getPassword()));
        user.setRole("USER");

        userRepository.save(user);

        return jwtUtil.generateToken(user.getUsername());
    }

    public String login(LoginRequest request) {
        Optional<User> user = userRepository.findByUsername(request.getUsername());
        if (user.isEmpty()) {
            throw new IllegalArgumentException("User not found");
        }
        if (!passwordEncoder.matches(request.getPassword(), user.get().getPassword())) {
            throw new IllegalArgumentException("Password mismatch");
        }
        return jwtUtil.generateToken(user.get().getUsername());
    }
}