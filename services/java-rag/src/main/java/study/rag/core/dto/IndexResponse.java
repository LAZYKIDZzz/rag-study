package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public record IndexResponse(
    @JsonProperty("document_id") String documentId,
    @JsonProperty("chunk_count") int chunkCount,
    @JsonProperty("vector_count") int vectorCount) {
}
