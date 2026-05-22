package study.rag.core;

import org.springframework.stereotype.Repository;
import study.rag.core.model.ChatMessage;
import study.rag.core.model.ChatSession;

import java.time.Instant;
import java.util.HashMap;
import java.util.Map;

@Repository
public class ChatRepository {
  private final Map<String, ChatSession> sessions = new HashMap<>();

  public synchronized ChatSession create(String userId) {
    ChatSession session = new ChatSession(Ids.next("ses"), userId, Instant.now());
    sessions.put(session.id(), session);
    return session;
  }

  public synchronized ChatSession get(String sessionId) {
    ChatSession session = sessions.get(sessionId);
    if (session == null) {
      throw new NotFoundException("Chat session " + sessionId + " was not found");
    }
    return session;
  }

  public synchronized ChatSession addMessage(String sessionId, ChatMessage message) {
    ChatSession session = get(sessionId);
    session.addMessage(message);
    return session;
  }
}
