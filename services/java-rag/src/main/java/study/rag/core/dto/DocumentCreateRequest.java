package study.rag.core.dto;

import jakarta.validation.constraints.NotBlank;

import java.util.Map;

public record DocumentCreateRequest(
    @NotBlank String title,
    @NotBlank String content,
    Map<String, Object> metadata) {
}
