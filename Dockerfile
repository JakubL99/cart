FROM alpine
RUN apk --no-cache add curl
ADD  cart /cart
ENTRYPOINT [ "/cart" ]