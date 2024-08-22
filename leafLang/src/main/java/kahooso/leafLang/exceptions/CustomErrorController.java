package kahooso.leafLang.exceptions;

import lombok.RequiredArgsConstructor;
import org.springframework.boot.web.error.ErrorAttributeOptions;
import org.springframework.boot.web.servlet.error.ErrorAttributes;
import org.springframework.boot.web.servlet.error.ErrorController;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.context.request.WebRequest;

import java.util.Map;

@Controller
@RequiredArgsConstructor
public class CustomErrorController implements ErrorController {

    private static final String PATH = "/error";

    private final ErrorAttributes errorAttributes;

    @RequestMapping(CustomErrorController.PATH)
    public ResponseEntity<ErrorDto> error(WebRequest webRequest) {

        Map<String, Object> attributes = errorAttributes.getErrorAttributes(
                webRequest,
                ErrorAttributeOptions.of(ErrorAttributeOptions.Include.EXCEPTION, ErrorAttributeOptions.Include.MESSAGE)
        );

        Integer status = (Integer) attributes.getOrDefault("status", HttpStatus.INTERNAL_SERVER_ERROR.value());
        String error = (String) attributes.getOrDefault("error", "Unknown error");
        String message = (String) attributes.getOrDefault("message", "No description");

        return ResponseEntity
                .status(status)
                .body(ErrorDto
                        .builder()
                        .error(error)
                        .errorDescription(message)
                        .build()
                );
    }
}
