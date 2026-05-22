package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Map;

public record RetrievedChunk(
    @JsonProperty("chunk_id") String chunkId,
    @JsonProperty("document_id") String documentId,
    @JsonProperty("document_title") String documentTitle,
    int ordinal,
    String content,
    double score,
    Map<String, Object> metadata) {
}
