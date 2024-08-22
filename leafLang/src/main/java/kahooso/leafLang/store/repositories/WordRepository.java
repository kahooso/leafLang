package kahooso.leafLang.store.repositories;

import kahooso.leafLang.store.entities.WordEntity;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;
import java.util.stream.Stream;

public interface WordRepository extends JpaRepository<WordEntity, Long> {

    Optional<WordEntity> findByInitial(String initial);

    Stream<WordEntity> streamAllByInitialStartsWithIgnoreCase(String initial);
    Stream<WordEntity> streamAllBy();
}
