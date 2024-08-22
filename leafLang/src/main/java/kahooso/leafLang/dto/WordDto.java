package kahooso.leafLang.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import kahooso.leafLang.Status;
import lombok.*;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class WordDto {

    @NonNull
    private long id;

    @NonNull
    private String initial;

    @NonNull
    private String translate;

    @NonNull
    private Status status;

    @NonNull
    @JsonProperty("example_of_usage")
    private String example_of_usage;
}
