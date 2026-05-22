package study.rag.core.model;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

public class ChatSession {
  private final String id;
  private final String userId;
  private final Instant createdAt;
  private final List<ChatMessage> messages = new ArrayList<>();

  public ChatSession(String id, String userId, Instant createdAt) {
    this.id = id;
    this.userId = userId;
    this.createdAt = createdAt;
  }

  public String id() {
    return id;
  }

  public String userId() {
    return userId;
  }

  public Instant createdAt() {
    return createdAt;
  }

  public List<ChatMessage> messages() {
    return List.copyOf(messages);
  }

  public void addMessage(ChatMessage message) {
    messages.add(message);
  }
}
