package study.rag.core.dto;

import java.util.Map;

public record ErrorBody(String code, String message, Map<String, Object> details) {
}
