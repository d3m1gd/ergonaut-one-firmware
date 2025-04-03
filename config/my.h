#ifndef My
#define My

#define __my_hash #
#define __my_pass(x) x
#define __pounded(a) __my_pass(__my_hash)a

#define gen_t0r(n)                      \
 t0r##n: t0r##n {                       \
     compatible = "zmk,behavior-macro"; \
     __pounded(binding-cells) = <0>;    \
     bindings = <&to 0 r##n>;           \
     label = "t0r";                     \
     wait-ms = <0>;                     \
 };

#define r32 &kpkp RG(M) M
#define r33 &kpkp RG(COMMA) COMMA
#define r34 &kpkp RG(DOT) DOT

#define MyMacros \
    gen_t0r(32) \
    gen_t0r(33) \
    gen_t0r(34)

#define MyKyes \
    X( 1, l11, a, b) \
    X( 2, l12, a, b) \
    X( 3, l13, a, b) \
    X( 4, l14, a, b) \
    X( 5, l15, a, b) \
    X( 6, l16, a, b) \
    X( 7, r11, a, b) \
    X( 8, r12, a, b) \
    X( 9, r13, a, b) \
    X(10, r14, a, b) \
    X(11, r15, a, b) \
    X(12, r16, a, b) \
    X(13, l21, a, b) \
    X(14, l22, a, b) \
    X(15, l23, a, b) \
    X(16, l24, a, b) \
    X(17, l25, a, b) \
    X(18, l26, a, b) \
    X(19, r21, a, b) \
    X(20, r22, a, b) \
    X(21, r23, a, b) \
    X(22, r24, a, b) \
    X(23, r25, a, b) \
    X(24, r26, a, b) \
    X(25, l31, a, b) \
    X(26, l32, a, b) \
    X(27, l33, a, b) \
    X(28, l34, a, b) \
    X(29, l35, a, b) \
    X(30, l36, a, b) \
    X(31, r31, a, b) \
    X(32, r32, a, b) \
    X(33, r33, a, b) \
    X(34, r34, a, b) \
    X(35, r35, a, b) \
    X(36, r36, a, b) \
    X(37, l41, a, b) \
    X(38, l42, a, b) \
    X(39, l43, a, b) \
    X(40, r41, a, b) \
    X(41, r42, a, b) \
    X(42, r43, a, b)

/*#define MyMoveLayer*/
/**/
/*&none &none &none &none &none &none &none    &none          &none        &none             &none &none*/
/*&none &none &none &none &none &none &kp LEFT &rmt LALT DOWN &rmt LGUI UP &rmt LSHIFT RIGHT &none &none*/
/*&none &none &none &none &none &none &none    &t0r32         &t0r33       &t0r34            &none &none*/
/*                  &none &none &to 0 &none    &none          &none                                I*/
/**/



#endif // !My
