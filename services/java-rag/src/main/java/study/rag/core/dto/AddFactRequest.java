package study.rag.core.dto;

import jakarta.validation.constraints.NotBlank;

public record AddFactRequest(@NotBlank String fact) {
}
