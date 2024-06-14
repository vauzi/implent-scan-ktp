package server

import (
	"strings"
)

func getValue(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func extractKTPData(text string) KTPData {
	lines := strings.Split(text, "\n")
	ktpData := KTPData{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "NIK"):
			ktpData.NIK = getValue(line)
		case strings.HasPrefix(line, "Nama"):
			ktpData.Nama = getValue(line)
		case strings.HasPrefix(line, "Tempat/Tgl Lahir"), strings.HasPrefix(line, "Tempat"):
			parts := strings.Split(getValue(line), " ")
			ktpData.TempatLahir = parts[0]
			if len(parts) > 1 {
				ktpData.TanggalLahir = parts[1]
			}
		case strings.HasPrefix(line, "Jenis Kelamin"):
			parts := strings.Split(getValue(line), " ")
			ktpData.JenisKelamin = parts[0]
			if len(parts) > 2 {
				ktpData.Alamat = parts[2]
			}
		case strings.HasPrefix(line, "Alamat"):
			ktpData.Alamat = getValue(line)
		case strings.HasPrefix(line, "RT/RW"), strings.HasPrefix(line, "RT"):
			ktpData.RT_RW = getValue(line)
		case strings.HasPrefix(line, "Kel/Desa"), strings.HasPrefix(line, "KellDesa"), strings.HasPrefix(line, "KelIDesa"), strings.HasPrefix(line, "Desa"):
			ktpData.KelDesa = getValue(line)
		case strings.HasPrefix(line, "Kecamatan"):
			ktpData.Kecamatan = getValue(line)
		case strings.HasPrefix(line, "Agama"):
			ktpData.Agama = getValue(line)
		case strings.HasPrefix(line, "Status Perkawinan"):
			ktpData.StatusPerkawinan = getValue(line)
		case strings.HasPrefix(line, "Pekerjaan"):
			ktpData.Pekerjaan = getValue(line)
		case strings.HasPrefix(line, "Kewarganegaraan"):
			ktpData.Kewarganegaraan = getValue(line)
		case strings.HasPrefix(line, "Berlaku Hingga"):
			ktpData.BerlakuHingga = getValue(line)
		}
	}

	return ktpData
}
