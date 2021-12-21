package main

import (
	"math"
	"math/cmplx"

	"github.com/fogleman/ln/ln"
)

func main() {
	eye := ln.Vector{-2, -2, 8}
	center := ln.Vector{0.1, 0, 0}
	up := ln.Vector{0, 1, 0}
	scene := ln.Scene{}
	scene.Add(CalabiYau(5, math.Pi/4, 16, -1, 1))
	width := 11.0 * 200
	height := 14.0 * 200
	fovy := 30.0
	paths := scene.Render(eye, center, up, width, height, fovy, 0.1, 10, 0.01)
	paths.WriteToPNG("out.png", width, height)
}

func CalabiYau(n int, alpha float64, count int, rmin float64, rmax float64) *ln.Mesh {
	cos := math.Cos(alpha)
	sin := math.Sin(alpha)

	dr := (rmax - rmin) / float64(count-1)
	di := (0.5 * math.Pi) / float64(count-1)

	vertices := []ln.Vector{}
	for k0 := 0; k0 < n; k0++ {
		for k1 := 0; k1 < n; k1++ {
			// Real and imaginary indices.
			for ir := 0; ir < count; ir++ {
				// Real and imaginary values.
				r := rmin + float64(ir)*dr

				for ii := 0; ii < count; ii++ {
					i := float64(ii) * di

					z0 := Z0k(r, i, float64(n), float64(k0))
					z1 := Z1k(r, i, float64(n), float64(k1))

					vertices = append(vertices, ln.Vector{X: real(z0), Y: real(z1), Z: cos*imag(z0) + sin*imag(z1)})
				}
			}
		}
	}

	vertexCount := count
	patchVertexCount := int(math.Pow(float64(vertexCount), 2))
	subdivisions := vertexCount - 1

	var triangles []*ln.Triangle

	for i := 0; i < n*n; i++ {
		offset := i * patchVertexCount

		for y := 0; y < subdivisions; y++ {
			for x := 0; x < subdivisions; x++ {
				v0 := y*vertexCount + x
				v1 := y*vertexCount + (x + 1)
				v2 := (y+1)*vertexCount + x
				v3 := (y+1)*vertexCount + (x + 1)

				v0 += offset
				v1 += offset
				v2 += offset
				v3 += offset

				/**
				 *   0     1
				 *    o---o
				 *    | \ |
				 *    o---o
				 *   2     3
				 */
				triangles = append(triangles, ln.NewTriangle(vertices[v0], vertices[v2], vertices[v3]))
				triangles = append(triangles, ln.NewTriangle(vertices[v0], vertices[v3], vertices[v1]))
			}
		}
	}

	return ln.NewMesh(triangles)
}

func PhaseFactor(k float64, n float64) complex128 {
	x := 2 * math.Pi * k / n
	return complex(math.Cos(x), math.Sin(x))
}

func U0(x complex128) complex128 {
	a := cmplx.Exp(x)
	b := cmplx.Exp(-x)
	return 0.5 * (a + b)
}

func U1(x complex128) complex128 {
	a := cmplx.Exp(x)
	b := cmplx.Exp(-x)
	return 0.5 * (a - b)
}

func Z0k(r float64, i float64, n float64, k float64) complex128 {
	phase := PhaseFactor(k, n)

	cos := U0(complex(r, i))
	powcos := cmplx.Pow(cos, complex(2/n, 0))

	return phase * powcos
}

func Z1k(r float64, i float64, n float64, k float64) complex128 {
	phase := PhaseFactor(k, n)

	sin := U1(complex(r, i))
	powsin := cmplx.Pow(sin, complex(2/n, 0))

	return phase * powsin
}
