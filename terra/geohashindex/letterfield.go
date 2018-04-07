package geohashindex

type LetterField struct {
	Amount int                    `firestore:"𝝨,omitempty"`
	Points map[string]interface{} `firestore:"Ω,omitempty"`
	A      *LetterField           `firestore:"a,omitempty"`
	B      *LetterField           `firestore:"b,omitempty"`
	C      *LetterField           `firestore:"c,omitempty"`
	D      *LetterField           `firestore:"d,omitempty"`
	E      *LetterField           `firestore:"e,omitempty"`
	F      *LetterField           `firestore:"f,omitempty"`
	G      *LetterField           `firestore:"g,omitempty"`
	H      *LetterField           `firestore:"h,omitempty"`
	I      *LetterField           `firestore:"i,omitempty"`
	J      *LetterField           `firestore:"j,omitempty"`
	K      *LetterField           `firestore:"k,omitempty"`
	L      *LetterField           `firestore:"l,omitempty"`
	M      *LetterField           `firestore:"m,omitempty"`
	N      *LetterField           `firestore:"n,omitempty"`
	O      *LetterField           `firestore:"o,omitempty"`
	P      *LetterField           `firestore:"p,omitempty"`
	Q      *LetterField           `firestore:"q,omitempty"`
	R      *LetterField           `firestore:"r,omitempty"`
	S      *LetterField           `firestore:"s,omitempty"`
	T      *LetterField           `firestore:"t,omitempty"`
	U      *LetterField           `firestore:"u,omitempty"`
	V      *LetterField           `firestore:"v,omitempty"`
	W      *LetterField           `firestore:"w,omitempty"`
	X      *LetterField           `firestore:"x,omitempty"`
	Y      *LetterField           `firestore:"y,omitempty"`
	Z      *LetterField           `firestore:"z,omitempty"`
	N0     *LetterField           `firestore:"0,omitempty"`
	N1     *LetterField           `firestore:"1,omitempty"`
	N2     *LetterField           `firestore:"2,omitempty"`
	N3     *LetterField           `firestore:"3,omitempty"`
	N4     *LetterField           `firestore:"4,omitempty"`
	N5     *LetterField           `firestore:"5,omitempty"`
	N6     *LetterField           `firestore:"6,omitempty"`
	N7     *LetterField           `firestore:"7,omitempty"`
	N8     *LetterField           `firestore:"8,omitempty"`
	N9     *LetterField           `firestore:"9,omitempty"`
}

func (Ω *LetterField) Add(hash string, i int, value interface{}) {
	Ω.Amount += 1
	if i >= 5 {
		if Ω.Points == nil {
			Ω.Points = map[string]interface{}{}
		}
		Ω.Points[hash] = value
		return
	}
	//Ω.Is = true
	switch hash[i] {
	case 'a':
		if Ω.A == nil {
			Ω.A = &LetterField{}
		}
		Ω.A.Add(hash, i+1, value)
	case 'b':
		if Ω.B == nil {
			Ω.B = &LetterField{}
		}
		Ω.B.Add(hash, i+1, value)
	case 'c':
		if Ω.C == nil {
			Ω.C = &LetterField{}
		}
		Ω.C.Add(hash, i+1, value)
	case 'd':
		if Ω.D == nil {
			Ω.D = &LetterField{}
		}
		Ω.D.Add(hash, i+1, value)
	case 'e':
		if Ω.E == nil {
			Ω.E = &LetterField{}
		}
		Ω.E.Add(hash, i+1, value)
	case 'f':
		if Ω.F == nil {
			Ω.F = &LetterField{}
		}
		Ω.F.Add(hash, i+1, value)
	case 'g':
		if Ω.G == nil {
			Ω.G = &LetterField{}
		}
		Ω.G.Add(hash, i+1, value)
	case 'h':
		if Ω.H == nil {
			Ω.H = &LetterField{}
		}
		Ω.H.Add(hash, i+1, value)
	case 'i':
		if Ω.I == nil {
			Ω.I = &LetterField{}
		}
		Ω.I.Add(hash, i+1, value)
	case 'j':
		if Ω.J == nil {
			Ω.J = &LetterField{}
		}
		Ω.J.Add(hash, i+1, value)
	case 'k':
		if Ω.K == nil {
			Ω.K = &LetterField{}
		}
		Ω.K.Add(hash, i+1, value)
	case 'l':
		if Ω.L == nil {
			Ω.L = &LetterField{}
		}
		Ω.L.Add(hash, i+1, value)
	case 'm':
		if Ω.M == nil {
			Ω.M = &LetterField{}
		}
		Ω.M.Add(hash, i+1, value)
	case 'n':
		if Ω.N == nil {
			Ω.N = &LetterField{}
		}
		Ω.N.Add(hash, i+1, value)
	case 'o':
		if Ω.O == nil {
			Ω.O = &LetterField{}
		}
		Ω.O.Add(hash, i+1, value)
	case 'p':
		if Ω.P == nil {
			Ω.P = &LetterField{}
		}
		Ω.P.Add(hash, i+1, value)
	case 'q':
		if Ω.Q == nil {
			Ω.Q = &LetterField{}
		}
		Ω.Q.Add(hash, i+1, value)
	case 'r':
		if Ω.R == nil {
			Ω.R = &LetterField{}
		}
		Ω.R.Add(hash, i+1, value)
	case 's':
		if Ω.S == nil {
			Ω.S = &LetterField{}
		}
		Ω.S.Add(hash, i+1, value)
	case 't':
		if Ω.T == nil {
			Ω.T = &LetterField{}
		}
		Ω.T.Add(hash, i+1, value)
	case 'u':
		if Ω.U == nil {
			Ω.U = &LetterField{}
		}
		Ω.U.Add(hash, i+1, value)
	case 'v':
		if Ω.V == nil {
			Ω.V = &LetterField{}
		}
		Ω.V.Add(hash, i+1, value)
	case 'w':
		if Ω.W == nil {
			Ω.W = &LetterField{}
		}
		Ω.W.Add(hash, i+1, value)
	case 'x':
		if Ω.X == nil {
			Ω.X = &LetterField{}
		}
		Ω.X.Add(hash, i+1, value)
	case 'y':
		if Ω.Y == nil {
			Ω.Y = &LetterField{}
		}
		Ω.Y.Add(hash, i+1, value)
	case 'z':
		if Ω.Z == nil {
			Ω.Z = &LetterField{}
		}
		Ω.Z.Add(hash, i+1, value)
	case '0':
		if Ω.N0 == nil {
			Ω.N0 = &LetterField{}
		}
		Ω.N0.Add(hash, i+1, value)
	case '1':
		if Ω.N1 == nil {
			Ω.N1 = &LetterField{}
		}
		Ω.N1.Add(hash, i+1, value)
	case '2':
		if Ω.N2 == nil {
			Ω.N2 = &LetterField{}
		}
		Ω.N2.Add(hash, i+1, value)
	case '3':
		if Ω.N3 == nil {
			Ω.N3 = &LetterField{}
		}
		Ω.N3.Add(hash, i+1, value)
	case '4':
		if Ω.N4 == nil {
			Ω.N4 = &LetterField{}
		}
		Ω.N4.Add(hash, i+1, value)
	case '5':
		if Ω.N5 == nil {
			Ω.N5 = &LetterField{}
		}
		Ω.N5.Add(hash, i+1, value)
	case '6':
		if Ω.N6 == nil {
			Ω.N6 = &LetterField{}
		}
		Ω.N6.Add(hash, i+1, value)
	case '7':
		if Ω.N7 == nil {
			Ω.N7 = &LetterField{}
		}
		Ω.N7.Add(hash, i+1, value)
	case '8':
		if Ω.N8 == nil {
			Ω.N8 = &LetterField{}
		}
		Ω.N8.Add(hash, i+1, value)
	case '9':
		if Ω.N9 == nil {
			Ω.N9 = &LetterField{}
		}
		Ω.N9.Add(hash, i+1, value)
	}
}
