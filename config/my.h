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

#define MyMacros \
    gen_t0r(32) \
    gen_t0r(33)

#endif // !My
