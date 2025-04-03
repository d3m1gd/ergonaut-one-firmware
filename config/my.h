#ifndef My
#define My

#define MyHash #
#define MyPass(x) x
#define MyPound(a) MyPass(MyHash)a

#define gen_macro_toand(LVL, RC, BEH)          \
 to##LVL##RC: to##LVL##RC {                    \
     compatible = "zmk,behavior-macro";        \
     MyPound(binding-cells) = <0>;             \
     bindings = <&to LVL BEH>;                 \
     label = MyCat(GenToAnd, RC);              \
     wait-ms = <0>;                            \
 };

#define MyStr(a) #a
#define MyCat(a, b) MyStr(a##b)
#define MyKp(a, b) &kp a
#define MyKpKp(a, b) &kpkp b a

#define MyKyes(X) \
    X(32, r32, M,     RG(M),     MyKpKp) \
    X(33, r33, COMMA, RG(COMMA), MyKpKp) \
    X(34, r34, DOT,   RG(DOT),   MyKpKp) \

    /*X( 1, l11, a, b, MyKp) \*/
    /*X( 2, l12, a, b, MyKp) \*/
    /*X( 3, l13, a, b, MyKp) \*/
    /*X( 4, l14, a, b, MyKp) \*/
    /*X( 5, l15, a, b, MyKp) \*/
    /*X( 6, l16, a, b, MyKp) \*/
    /*X( 7, r11, a, b, MyKp) \*/
    /*X( 8, r12, a, b, MyKp) \*/
    /*X( 9, r13, a, b, MyKp) \*/
    /*X(10, r14, a, b, MyKp) \*/
    /*X(11, r15, a, b, MyKp) \*/
    /*X(12, r16, a, b, MyKp) \*/
    /*X(13, l21, a, b, MyKp) \*/
    /*X(14, l22, a, b, MyKp) \*/
    /*X(15, l23, a, b, MyKp) \*/
    /*X(16, l24, a, b, MyKp) \*/
    /*X(17, l25, a, b, MyKp) \*/
    /*X(18, l26, a, b, MyKp) \*/
    /*X(19, r21, a, b, MyKp) \*/
    /*X(20, r22, a, b, MyKp) \*/
    /*X(21, r23, a, b, MyKp) \*/
    /*X(22, r24, a, b, MyKp) \*/
    /*X(23, r25, a, b, MyKp) \*/
    /*X(24, r26, a, b, MyKp) \*/
    /*X(25, l31, a, b, MyKp) \*/
    /*X(26, l32, a, b, MyKp) \*/
    /*X(27, l33, a, b, MyKp) \*/
    /*X(28, l34, a, b, MyKp) \*/
    /*X(29, l35, a, b, MyKp) \*/
    /*X(30, l36, a, b, MyKp) \*/
    /*X(31, r31, a, b, MyKp) \*/
    /**/
    /*X(35, r35, a, b, MyKp) \*/
    /*X(36, r36, a, b, MyKp) \*/
    /*X(37, l41, a, b, MyKp) \*/
    /*X(38, l42, a, b, MyKp) \*/
    /*X(39, l43, a, b, MyKp) \*/
    /*X(40, r41, a, b, MyKp) \*/
    /*X(41, r42, a, b, MyKp) \*/
    /*X(42, r43, a, b, MyKp)*/

#define MyGenTo0(num, rowcol, key, alt, con) gen_macro_toand(0, rowcol, con(key, alt))

/*#define MyMoveLayer*/
/**/
/*&none &none &none &none &none &none &none    &none          &none        &none             &none &none*/
/*&none &none &none &none &none &none &kp LEFT &rmt LALT DOWN &rmt LGUI UP &rmt LSHIFT RIGHT &none &none*/
/*&none &none &none &none &none &none &none    &t0r32         &t0r33       &t0r34            &none &none*/
/*                  &none &none &to 0 &none    &none          &none                                I*/
/**/



#endif // !My
