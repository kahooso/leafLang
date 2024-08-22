package kahooso.leafLang.store.entities;

import jakarta.persistence.*;
import kahooso.leafLang.Status;
import lombok.*;

@Getter
@Setter
@Entity
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table(name = "word")
public class WordEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private long id;

    private String initial;

    private String translate;

    private Status status;

    private String example_of_usage;
}
