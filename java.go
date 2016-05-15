package colfer

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateJava writes the code into the respective ".java" files.
func GenerateJava(basedir string, packages []*Package) error {
	t := template.New("java-code").Delims("<:", ":>")
	template.Must(t.Parse(javaCode))

	for _, p := range packages {
		p.NameNative = strings.Replace(p.Name, "/", ".", -1)
	}

	for _, p := range packages {
		pkgdir, err := makePkgDir(p, basedir)
		if err != nil {
			return err
		}

		for _, s := range p.Structs {
			for _, f := range s.Fields {
				switch f.Type {
				default:
					if f.TypeRef == nil {
						f.TypeNative = f.Type
					} else {
						f.TypeNative = f.TypeRef.NameTitle()
						if f.TypeRef.Pkg != p {
							f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
						}
					}
				case "bool":
					f.TypeNative = "boolean"
				case "uint32", "int32":
					f.TypeNative = "int"
				case "uint64", "int64":
					f.TypeNative = "long"
				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"
				case "timestamp":
					f.TypeNative = "java.time.Instant"
				case "text":
					f.TypeNative = "String"
				case "binary":
					f.TypeNative = "byte[]"
				}
			}

			f, err := os.Create(filepath.Join(pkgdir, s.NameTitle()+".java"))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := t.Execute(f, s); err != nil {
				return err
			}
		}
	}
	return nil
}

const javaCode = `package <:.Pkg.NameNative:>;


// This file was generated by colf(1); DO NOT EDIT


import static java.lang.String.format;
import java.nio.BufferOverflowException;
import java.nio.BufferUnderflowException;
import javax.xml.bind.TypeConstraintException;
import javax.xml.bind.DataBindingException;


/**
 * Data bean with built-in serialization support.
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class <:.NameTitle:> implements java.io.Serializable {

	/** The upper limit for serial byte sizes. */
	public static int colferSizeMax = 16 * 1024 * 1024;

	/** The upper limit for the number of elements in a list. */
	public static int colferListMax = 64 * 1024;

	private static final java.nio.charset.Charset _utf8 = java.nio.charset.Charset.forName("UTF-8");
<:- range .Fields:>
<:- if eq .Type "binary":>
	private static final byte[] _zero<:.NameTitle:> = new byte[0];
<:- else if .TypeArray:>
	private static final <:.TypeNative:>[] _zero<:.NameTitle:> = new <:.TypeNative:>[0];
<:- end:>
<:- end:>
<:range .Fields:>
	public <:.TypeNative:><:if .TypeArray:>[]<:end:> <:.Name:>
<:- if eq .Type "text":> = ""
<:- else if eq .Type "binary":> = _zero<:.NameTitle:>
<:- else if .TypeArray:> = _zero<:.NameTitle:>
<:- end:>;<:end:>


	/**
	 * Serializes the object.
<:- range .Fields:><:if .TypeArray:>
	 * All {@code null} entries in {@link #<:.Name:>} will be replaced with a {@code new} value.
<:- end:><:end:>
	 * @param buf the data destination.
	 * @param offset the first byte index.
	 * @return the index of the first byte after the last byte written.
	 * @throws BufferOverflowException when {@code buf} is too small.
	 * @throws IllegalStateException on an upper limit breach defined by either {@link #colferSizeMax} or {@link #colferListMax}.
	 */
	public int marshal(byte[] buf, int offset) {
		int i = offset;
		try {
<:range .Fields:><:if eq .Type "bool":>
			if (this.<:.Name:>) {
				buf[i++] = (byte) <:.Index:>;
			}
<:else if eq .Type "uint32":>
			if (this.<:.Name:> != 0) {
				int x = this.<:.Name:>;
				if ((x & ~((1 << 21) - 1)) != 0) {
					buf[i++] = (byte) (<:.Index:> | 0x80);
					buf[i++] = (byte) (x >>> 24);
					buf[i++] = (byte) (x >>> 16);
					buf[i++] = (byte) (x >>> 8);
					buf[i++] = (byte) (x);
				} else {
					buf[i++] = (byte) <:.Index:>;
					while ((x & ~((1 << 7) - 1)) != 0) {
						buf[i++] = (byte) (x | 0x80);
						x >>>= 7;
					}
					buf[i++] = (byte) x;
				}
			}
<:else if eq .Type "uint64":>
			if (this.<:.Name:> != 0) {
				long x = this.<:.Name:>;
				if ((x & ~((1 << 49) - 1)) != 0) {
					buf[i++] = (byte) (<:.Index:> | 0x80);
					buf[i++] = (byte) (x >>> 56);
					buf[i++] = (byte) (x >>> 48);
					buf[i++] = (byte) (x >>> 40);
					buf[i++] = (byte) (x >>> 32);
					buf[i++] = (byte) (x >>> 24);
					buf[i++] = (byte) (x >>> 16);
					buf[i++] = (byte) (x >>> 8);
					buf[i++] = (byte) (x);
				} else {
					buf[i++] = (byte) <:.Index:>;
					for (int n = 0; n < 8 && (x & ~((1L << 7) - 1)) != 0; n++) {
						buf[i++] = (byte) (x | 0x80);
						x >>>= 7;
					}
					buf[i++] = (byte) x;
				}
			}
<:else if eq .Type "int32":>
			if (this.<:.Name:> != 0) {
				int x = this.<:.Name:>;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (<:.Index:> | 0x80);
				} else
					buf[i++] = (byte) <:.Index:>;
				while ((x & ~((1 << 7) - 1)) != 0) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}
<:else if eq .Type "int64":>
			if (this.<:.Name:> != 0) {
				long x = this.<:.Name:>;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (<:.Index:> | 0x80);
				} else
					buf[i++] = (byte) <:.Index:>;
				for (int n = 0; n < 8 && (x & ~((1L << 7) - 1)) != 0; n++) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}
<:else if eq .Type "float32":>
			if (this.<:.Name:> != 0.0f) {
				buf[i++] = (byte) <:.Index:>;
				int x = Float.floatToRawIntBits(this.<:.Name:>);
				buf[i++] = (byte) (x >>> 24);
				buf[i++] = (byte) (x >>> 16);
				buf[i++] = (byte) (x >>> 8);
				buf[i++] = (byte) (x);
			}
<:else if eq .Type "float64":>
			if (this.<:.Name:> != 0.0) {
				buf[i++] = (byte) <:.Index:>;
				long x = Double.doubleToRawLongBits(this.<:.Name:>);
				buf[i++] = (byte) (x >>> 56);
				buf[i++] = (byte) (x >>> 48);
				buf[i++] = (byte) (x >>> 40);
				buf[i++] = (byte) (x >>> 32);
				buf[i++] = (byte) (x >>> 24);
				buf[i++] = (byte) (x >>> 16);
				buf[i++] = (byte) (x >>> 8);
				buf[i++] = (byte) (x);
			}
<:else if eq .Type "timestamp":>
			if (this.<:.Name:> != null) {
				long s = this.<:.Name:>.getEpochSecond();
				int ns = this.<:.Name:>.getNano();
				if (s != 0 || ns != 0) {
					if (s >= 0 && s < (1L << 32)) {
						buf[i++] = (byte) <:.Index:>;
						buf[i++] = (byte) (s >>> 24);
						buf[i++] = (byte) (s >>> 16);
						buf[i++] = (byte) (s >>> 8);
						buf[i++] = (byte) (s);
						buf[i++] = (byte) (ns >>> 24);
						buf[i++] = (byte) (ns >>> 16);
						buf[i++] = (byte) (ns >>> 8);
						buf[i++] = (byte) (ns);
					} else {
						buf[i++] = (byte) (<:.Index:> | 0x80);
						buf[i++] = (byte) (s >>> 56);
						buf[i++] = (byte) (s >>> 48);
						buf[i++] = (byte) (s >>> 40);
						buf[i++] = (byte) (s >>> 32);
						buf[i++] = (byte) (s >>> 24);
						buf[i++] = (byte) (s >>> 16);
						buf[i++] = (byte) (s >>> 8);
						buf[i++] = (byte) (s);
						buf[i++] = (byte) (ns >>> 24);
						buf[i++] = (byte) (ns >>> 16);
						buf[i++] = (byte) (ns >>> 8);
						buf[i++] = (byte) (ns);
					}
				}
			}
<:else if eq .Type "text":>
			if (! this.<:.Name:>.isEmpty()) {
				buf[i++] = (byte) <:.Index:>;
				String s = this.<:.Name:>;
				int sLength = s.length();

				int start = ++i;
				for (int sIndex = 0; sIndex < sLength; sIndex++) {
					char c = s.charAt(sIndex);
					if (c < 128) {
						buf[i++] = (byte) c;
					} else if (c < 2048) {
						buf[i++] = (byte) (192 | c >>> 6);
						buf[i++] = (byte) (128 | c & 63);
					} else if (! Character.isSurrogate(c)) {
						buf[i++] = (byte) (224 | c >>> 12);
						buf[i++] = (byte) (128 | c >>> 6 & 63);
						buf[i++] = (byte) (128 | c & 63);
					} else if (++sIndex != sLength) {
						int cp = Character.toCodePoint(c, s.charAt(sIndex));
						buf[i++] = (byte) (240 | cp >>> 18);
						buf[i++] = (byte) (128 | cp >>> 12 & 63);
						buf[i++] = (byte) (128 | cp >>> 6 & 63);
						buf[i++] = (byte) (128 | cp & 63);
					}
				}

				int size = i - start;
				if (size > colferSizeMax)
					throw new IllegalStateException(format("colfer: field <:.String:> size %d exceeds %d UTF-8 bytes", size, colferSizeMax));

				int shift = 0;
				for (int x = size; (x & ~((1 << 7) - 1)) != 0; x >>>= 7) shift++;
				if (shift != 0) System.arraycopy(buf, start, buf, start + shift, size);
				i = start + shift + size;

				start--;
				while ((size & ~((1 << 7) - 1)) != 0) {
					buf[start++] = (byte) (size | 0x80);
					size >>>= 7;
				}
				buf[start++] = (byte) size;
			}
<:else if eq .Type "binary":>
			if (this.<:.Name:>.length != 0) {
				buf[i++] = (byte) <:.Index:>;

				int x = this.<:.Name:>.length;
				if (x > colferSizeMax)
					throw new IllegalStateException(format("colfer: field <:.String:> size %d exceeds %d bytes", x, colferSizeMax));
				while ((x & ~((1 << 7) - 1)) != 0) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;

				System.arraycopy(this.<:.Name:>, 0, buf, i, this.<:.Name:>.length);
				i += this.<:.Name:>.length;
			}
<:else if .TypeArray:>
			if (this.<:.Name:>.length != 0) {
				buf[i++] = (byte) <:.Index:>;
				<:.TypeNative:>[] a = this.<:.Name:>;

				int x = a.length;
				if (x > colferListMax)
					throw new IllegalStateException(format("colfer: field <:.String:> length %d exceeds %d elements", x, colferListMax));
				while ((x & ~((1 << 7) - 1)) != 0) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;

				for (int ai = 0; ai < a.length; ai++) {
					<:.TypeNative:> o = a[ai];
					if (o == null) {
						o = new <:.TypeNative:>();
						a[ai] = o;
					}
					i = o.marshal(buf, i);
				}
			}
<:else:>
			if (this.<:.Name:> != null) {
				buf[i++] = (byte) <:.Index:>;
				i = this.<:.Name:>.marshal(buf, i);
			}
<:end:><:end:>
			buf[i++] = (byte) 0x7f;
			return i;
		} catch (IndexOutOfBoundsException e) {
			if (i - offset > colferSizeMax)
				throw new IllegalStateException(format("colfer: serial exceeds %d bytes", colferSizeMax));
			if (i >= buf.length)
				throw new BufferOverflowException();
			throw new RuntimeException("colfer: bug", e);
		}
	}

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the first byte index.
	 * @return the index of the first byte after the last byte read.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws TypeConstraintException on an upper limit breach defined by either {@link #colferSizeMax} or {@link #colferListMax}.
	 * @throws DataBindingException when the data does not match this object's schema.
	 */
	public int unmarshal(byte[] buf, int offset)
	throws BufferUnderflowException, TypeConstraintException, DataBindingException {
		int i = offset;
		try {
			byte header = buf[i++];
<:range .Fields:><:if eq .Type "bool":>
			if (header == (byte) <:.Index:>) {
				this.<:.Name:> = true;
				header = buf[i++];
			}
<:else if eq .Type "uint32":>
			if (header == (byte) <:.Index:>) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				this.<:.Name:> = (buf[i++] & 0xff) << 24 | (buf[i++] & 0xff) << 16 | (buf[i++] & 0xff) << 8 | (buf[i++] & 0xff);
				header = buf[i++];
			}
<:else if eq .Type "uint64":>
			if (header == (byte) <:.Index:>) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				this.<:.Name:> = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				header = buf[i++];
			}
<:else if eq .Type "int32":>
			if (header == (byte) <:.Index:>) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.<:.Name:> = -x;
				header = buf[i++];
			}
<:else if eq .Type "int64":>
			if (header == (byte) <:.Index:>) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = x;
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.<:.Name:> = -x;
				header = buf[i++];
			}
<:else if eq .Type "float32":>
			if (header == (byte) <:.Index:>) {
				int x = (buf[i++] & 0xff) << 24 | (buf[i++] & 0xff) << 16 | (buf[i++] & 0xff) << 8 | (buf[i++] & 0xff);
				this.<:.Name:> = Float.intBitsToFloat(x);
				header = buf[i++];
			}
<:else if eq .Type "float64":>
			if (header == (byte) <:.Index:>) {
				long x = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = Double.longBitsToDouble(x);
				header = buf[i++];
			}
<:else if eq .Type "timestamp":>
			if (header == (byte) <:.Index:>) {
				long s = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				long ns = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = java.time.Instant.ofEpochSecond(s, ns);
				header = buf[i++];
			} else if (header == (byte) (<:.Index:> | 0x80)) {
				long s = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				long ns = (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.<:.Name:> = java.time.Instant.ofEpochSecond(s, ns);
				header = buf[i++];
			}
<:else if eq .Type "text":>
			if (header == (byte) <:.Index:>) {
				int n = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					n |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (n > colferSizeMax)
					throw new TypeConstraintException(format("colfer: field <:.String:> size %d exceeds %d UTF-8 bytes", n, colferSizeMax));
				this.<:.Name:> = new String(buf, i, n, this._utf8);
				i += n;
				header = buf[i++];
			}
<:else if eq .Type "binary":>
			if (header == (byte) <:.Index:>) {
				int n = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					n |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (n > colferSizeMax)
					throw new TypeConstraintException(format("colfer: field <:.String:> size %d exceeds %d bytes", n, colferSizeMax));
				this.<:.Name:> = new byte[n];
				System.arraycopy(buf, i, this.<:.Name:>, 0, n);
				i += n;
				header = buf[i++];
			}
<:else if .TypeArray:>
			if (header == (byte) <:.Index:>) {
				int n = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					n |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (n > colferListMax)
					throw new TypeConstraintException(format("colfer: field <:.String:> length %d exceeds %d elements", n, colferListMax));
				<:.TypeNative:>[] a = new <:.TypeNative:>[n];
				for (int ai = 0; ai < n; ai++) {
					<:.TypeNative:> o = new <:.TypeNative:>();
					i = o.unmarshal(buf, i);
					a[ai] = o;
				}
				this.<:.Name:> = a;
				header = buf[i++];
			}
<:else:>
			if (header == (byte) <:.Index:>) {
				this.<:.Name:> = new <:.TypeNative:>();
				i = this.<:.Name:>.unmarshal(buf, i);
				header = buf[i++];
			}
<:end:><:end:>
			if (header != (byte) 0x7f)
				throw new DataBindingException(format("colfer: unknown header at byte %d", i - 1), null);
		} catch (IndexOutOfBoundsException e) {
			if (i - offset > colferSizeMax)
				throw new TypeConstraintException(format("colfer: serial exceeds %d bytes", colferSizeMax));
			if (i >= buf.length)
				throw new BufferUnderflowException();
			throw new RuntimeException("colfer: bug", e);
		}

		return i;
	}
<:range .Fields:>
	public <:.TypeNative:><:if .TypeArray:>[]<:end:> get<:.NameTitle:>() {
		return this.<:.Name:>;
	}

	public void set<:.NameTitle:>(<:.TypeNative:><:if .TypeArray:>[]<:end:> value) {
		this.<:.Name:> = value;
	}
<:end:>
	@Override
	public final int hashCode() {
		return java.util.Objects.hash(0x7f<:range .Fields:>, <:.Name:><:end:>);
	}

	@Override
	public final boolean equals(Object o) {
		return o instanceof <:.NameTitle:> && equals((<:.NameTitle:>) o);
	}

	public final boolean equals(<:.NameTitle:> o) {
		return o != null
<:- range .Fields:>
<:- if eq .Type "bool" "uint32" "uint64" "int32" "int64":>
			&& this.<:.Name:> == o.<:.Name:>
<:- else if eq .Type "binary":>
			&& java.util.Arrays.equals(this.<:.Name:>, o.<:.Name:>)
<:- else if .TypeArray:>
			&& java.util.Arrays.equals(this.<:.Name:>, o.<:.Name:>)
<:- else:>
			&& java.util.Objects.equals(this.<:.Name:>, o.<:.Name:>)
<:- end:><:end:>;
	}

}
`
