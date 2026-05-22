package study.rag.api;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import study.rag.core.RagService;
import study.rag.core.dto.AddFactRequest;
import study.rag.core.dto.CapabilitiesResponse;
import study.rag.core.dto.ChatMessageRequest;
import study.rag.core.dto.ChatMessageResponse;
import study.rag.core.dto.ChatSessionCreateRequest;
import study.rag.core.dto.ChatSessionResponse;
import study.rag.core.dto.DocumentCreateRequest;
import study.rag.core.dto.DocumentListResponse;
import study.rag.core.dto.DocumentResponse;
import study.rag.core.dto.HealthResponse;
import study.rag.core.dto.IndexResponse;
import study.rag.core.dto.MemoryResponse;
import study.rag.core.dto.RetrievalRequest;
import study.rag.core.dto.RetrievalResponse;

@CrossOrigin(origins = "*")
@RestController
public class RagController {
  private final RagService ragService;

  public RagController(RagService ragService) {
    this.ragService = ragService;
  }

  @GetMapping("/health")
  public HealthResponse health() {
    return new HealthResponse("ok", "java-rag", "0.1.0");
  }

  @GetMapping("/capabilities")
  public CapabilitiesResponse capabilities() {
    return ragService.capabilities();
  }

  @PostMapping("/documents")
  @ResponseStatus(HttpStatus.CREATED)
  public DocumentResponse createDocument(@Valid @RequestBody DocumentCreateRequest request) {
    return ragService.createDocument(request);
  }

  @GetMapping("/documents")
  public DocumentListResponse listDocuments() {
    return ragService.listDocuments();
  }

  @PostMapping("/documents/{documentId}/index")
  public IndexResponse indexDocument(@PathVariable String documentId) {
    return ragService.indexDocument(documentId);
  }

  @PostMapping("/retrieval/search")
  public RetrievalResponse search(@Valid @RequestBody RetrievalRequest request) {
    return ragService.search(request);
  }

  @PostMapping("/chat/sessions")
  @ResponseStatus(HttpStatus.CREATED)
  public ChatSessionResponse createSession(
      @RequestParam(defaultValue = "demo-user") String userId,
      @RequestBody(required = false) ChatSessionCreateRequest request) {
    String resolvedUserId = request != null && request.userId() != null && !request.userId().isBlank()
        ? request.userId()
        : userId;
    return ragService.createSession(resolvedUserId);
  }

  @PostMapping("/chat/sessions/{sessionId}/messages")
  public ChatMessageResponse sendMessage(
      @PathVariable String sessionId,
      @Valid @RequestBody ChatMessageRequest request) {
    return ragService.sendMessage(sessionId, request);
  }

  @GetMapping("/chat/sessions/{sessionId}")
  public ChatSessionResponse getSession(@PathVariable String sessionId) {
    return ragService.getSession(sessionId);
  }

  @GetMapping("/memory/users/{userId}")
  public MemoryResponse getMemory(@PathVariable String userId) {
    return ragService.getMemory(userId);
  }

  @PostMapping("/memory/users/{userId}/facts")
  @ResponseStatus(HttpStatus.CREATED)
  public MemoryResponse addFact(@PathVariable String userId, @Valid @RequestBody AddFactRequest request) {
    return ragService.addFact(userId, request.fact());
  }
}
