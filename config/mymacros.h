#ifndef MyMacros
#define MyMacros

#include "my.h"

#define MyGenTo0(num, rowcol, key, alt, con) gen_macro_toand(0, rowcol, con(key, alt))

MyKyes(MyGenTo0)

/*#define MyMoveLayer*/
/**/
/*&none &none &none &none &none &none &none    &none          &none        &none             &none &none*/
/*&none &none &none &none &none &none &kp LEFT &rmt LALT DOWN &rmt LGUI UP &rmt LSHIFT RIGHT &none &none*/
/*&none &none &none &none &none &none &none    &t0r32         &t0r33       &t0r34            &none &none*/
/*                  &none &none &to 0 &none    &none          &none                                I*/
/**/



#endif // !MyMacros
